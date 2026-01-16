package service

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/airsss993/work-svc/internal/config"
	"github.com/airsss993/work-svc/internal/domain"
	"github.com/airsss993/work-svc/pkg/logger"
	"github.com/go-ldap/ldap/v3"
)

type GroupService interface {
	GetITGroups(ctx context.Context) ([]domain.GroupInfo, error)
	GetGroupStudents(ctx context.Context, groupName string) ([]domain.Student, error)
	GetGroupSubgroups(ctx context.Context, groupName string) (*domain.SubgroupsResponse, error)
	GetGroupStudentsFiltered(ctx context.Context, groupName, subgroup string) ([]domain.Student, error)
}

type GroupServiceImpl struct {
	cfg    *config.Config
	appCfg *config.App
}

func NewGroupService(cfg *config.Config, appCfg *config.App) *GroupServiceImpl {
	return &GroupServiceImpl{
		cfg:    cfg,
		appCfg: appCfg,
	}
}

func (s *GroupServiceImpl) GetITGroups(ctx context.Context) ([]domain.GroupInfo, error) {
	if ctx.Err() != nil {
		return []domain.GroupInfo{}, nil
	}

	if s.appCfg.Test {
		return []domain.GroupInfo{
			{Name: "ИТ24-11"},
			{Name: "ИТ24-12"},
			{Name: "ИТ23-11"},
		}, nil
	}

	l, err := ldap.DialURL(s.cfg.LDAP.URL)
	if err != nil {
		logger.Error(fmt.Errorf("failed to connect to LDAP: %w", err))
		return nil, fmt.Errorf("LDAP connection failed")
	}
	defer l.Close()

	baseDN := "dc=it-college,dc=ru"
	searchRequest := ldap.NewSearchRequest(
		fmt.Sprintf("ou=Current,%s", baseDN),
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		"(|(objectClass=group)(objectClass=groupOfNames)(objectClass=posixGroup))",
		[]string{"cn", "dn"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		logger.Error(fmt.Errorf("LDAP search failed: %w", err))
		return nil, fmt.Errorf("search failed")
	}

	var groups []domain.GroupInfo
	for _, entry := range sr.Entries {
		cn := entry.GetAttributeValue("cn")

		if strings.HasPrefix(cn, "ИТ") && cn != "" {
			groups = append(groups, domain.GroupInfo{
				Name: cn,
			})
		}
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Name > groups[j].Name
	})

	return groups, nil
}

func (s *GroupServiceImpl) GetGroupStudents(ctx context.Context, groupName string) ([]domain.Student, error) {
	if ctx.Err() != nil {
		return []domain.Student{}, nil
	}

	if groupName == "" {
		return nil, fmt.Errorf("group name is required")
	}

	if s.appCfg.Test {
		return []domain.Student{
			{ID: "i24s0291", Username: "Коломацкий Иван", PhotoURL: ""},
			{ID: "i24s0002", Username: "Джапаридзе Артем", PhotoURL: "/api/photos/2024291.png"},
			{ID: "i24s0001", Username: "Тестов Тест", PhotoURL: ""},
		}, nil
	}

	l, err := ldap.DialURL(s.cfg.LDAP.URL)
	if err != nil {
		logger.Error(fmt.Errorf("failed to connect to LDAP: %w", err))
		return nil, fmt.Errorf("LDAP connection failed")
	}
	defer l.Close()

	baseDN := "dc=it-college,dc=ru"

	groupFilter := fmt.Sprintf("(&(|(objectClass=group)(objectClass=groupOfNames)(objectClass=posixGroup))(cn=%s))", ldap.EscapeFilter(groupName))
	groupSearchRequest := ldap.NewSearchRequest(
		fmt.Sprintf("ou=Current,%s", baseDN),
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		groupFilter,
		[]string{"member"},
		nil,
	)

	groupResult, err := l.Search(groupSearchRequest)
	if err != nil {
		logger.Error(fmt.Errorf("LDAP group search failed: %w", err))
		return nil, fmt.Errorf("group search failed")
	}

	if len(groupResult.Entries) == 0 {
		return []domain.Student{}, nil
	}

	groupEntry := groupResult.Entries[0]
	memberDNs := groupEntry.GetAttributeValues("member")

	if len(memberDNs) == 0 {
		return []domain.Student{}, nil
	}

	uidRegex := regexp.MustCompile(`uid=([^,]+)`)
	var uids []string
	for _, memberDN := range memberDNs {
		matches := uidRegex.FindStringSubmatch(memberDN)
		if len(matches) > 1 {
			uids = append(uids, matches[1])
		}
	}

	if len(uids) == 0 {
		return []domain.Student{}, nil
	}

	var filterParts []string
	for _, uid := range uids {
		filterParts = append(filterParts, fmt.Sprintf("(uid=%s)", ldap.EscapeFilter(uid)))
	}
	studentsFilter := fmt.Sprintf("(&(objectClass=person)(|%s))", strings.Join(filterParts, ""))

	studentsSearchRequest := ldap.NewSearchRequest(
		fmt.Sprintf("ou=people,%s", baseDN),
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		studentsFilter,
		[]string{"uid", "cn", "employeeNumber"},
		nil,
	)

	studentsResult, err := l.Search(studentsSearchRequest)
	if err != nil {
		logger.Error(fmt.Errorf("LDAP students search failed: %w", err))
		return nil, fmt.Errorf("students search failed")
	}

	var students []domain.Student
	for _, entry := range studentsResult.Entries {
		uid := entry.GetAttributeValue("uid")
		cn := entry.GetAttributeValue("cn")
		en := entry.GetAttributeValue("employeeNumber")

		if uid != "" && cn != "" && en != "" {
			var photoURL string
			photoPath := fmt.Sprintf("./photos/%s.png", en)

			if _, err := os.Stat(photoPath); err == nil {
				photoURL = fmt.Sprintf("/api/photos/%s.png", en)
			} else if en[0:2] == "23" {
				photoPathAlt := fmt.Sprintf("./photos/20%s%s.png", en[0:2], en[3:])
				if _, err := os.Stat(photoPathAlt); err == nil {
					photoURL = fmt.Sprintf("/api/photos/20%s%s.png", en[0:2], en[3:])
				}
			}

			students = append(students, domain.Student{
				ID:       uid,
				Username: cn,
				PhotoURL: photoURL,
			})
		}
	}

	return students, nil
}

func calculateCourse(groupName string) (int, error) {
	re := regexp.MustCompile(`ИТ(\d{2})-\d{2}`)
	matches := re.FindStringSubmatch(groupName)
	if len(matches) < 2 {
		return 0, fmt.Errorf("invalid group name format")
	}

	yearShort, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, fmt.Errorf("invalid year in group name")
	}
	year := 2000 + yearShort

	now := time.Now()
	academicYear := now.Year()
	if now.Month() < time.September {
		academicYear--
	}

	course := academicYear - year + 1
	return course, nil
}

func (s *GroupServiceImpl) getGroupMembers(l *ldap.Conn, groupCN string, baseDN string) ([]string, error) {
	groupFilter := fmt.Sprintf("(&(|(objectClass=group)(objectClass=groupOfNames)(objectClass=posixGroup))(cn=%s))", ldap.EscapeFilter(groupCN))
	groupSearchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		groupFilter,
		[]string{"member"},
		nil,
	)

	groupResult, err := l.Search(groupSearchRequest)
	if err != nil {
		return nil, fmt.Errorf("LDAP group search failed: %w", err)
	}

	if len(groupResult.Entries) == 0 {
		return []string{}, nil
	}

	return groupResult.Entries[0].GetAttributeValues("member"), nil
}

func hasIntersection(members []string, set map[string]bool) bool {
	for _, member := range members {
		if set[member] {
			return true
		}
	}
	return false
}

func membersToSet(members []string) map[string]bool {
	set := make(map[string]bool)
	for _, m := range members {
		set[m] = true
	}
	return set
}

func intersectMembers(members1, members2 []string) []string {
	set := membersToSet(members2)
	var result []string
	for _, m := range members1 {
		if set[m] {
			result = append(result, m)
		}
	}
	return result
}

func (s *GroupServiceImpl) GetGroupSubgroups(ctx context.Context, groupName string) (*domain.SubgroupsResponse, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if groupName == "" {
		return nil, fmt.Errorf("group name is required")
	}

	course, err := calculateCourse(groupName)
	if err != nil {
		return nil, err
	}

	if s.appCfg.Test {
		response := &domain.SubgroupsResponse{
			English: []string{"A0.21", "A1.21", "A1.22", "A2.21", "B1.21"},
			Course:  course,
		}
		if course >= 2 && course <= 3 {
			response.Profiles = []string{"BE", "FE", "CD", "GD", "PM", "SA"}
		}
		return response, nil
	}

	l, err := ldap.DialURL(s.cfg.LDAP.URL)
	if err != nil {
		logger.Error(fmt.Errorf("failed to connect to LDAP: %w", err))
		return nil, fmt.Errorf("LDAP connection failed")
	}
	defer l.Close()

	baseDN := "dc=it-college,dc=ru"
	currentBaseDN := fmt.Sprintf("ou=Current,%s", baseDN)

	mainGroupMembers, err := s.getGroupMembers(l, groupName, currentBaseDN)
	if err != nil {
		logger.Error(err)
		return nil, fmt.Errorf("failed to get main group members")
	}

	mainGroupSet := membersToSet(mainGroupMembers)

	courseDigit := strconv.Itoa(course)
	englishFilter := fmt.Sprintf("(&(|(objectClass=group)(objectClass=groupOfNames)(objectClass=posixGroup))(|(cn=A0.%s*)(cn=A1.%s*)(cn=A2.%s*)(cn=B1.%s*)))", courseDigit, courseDigit, courseDigit, courseDigit)

	englishSearchRequest := ldap.NewSearchRequest(
		currentBaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		englishFilter,
		[]string{"cn", "member"},
		nil,
	)

	englishResult, err := l.Search(englishSearchRequest)
	if err != nil {
		logger.Error(fmt.Errorf("LDAP english subgroups search failed: %w", err))
		return nil, fmt.Errorf("english subgroups search failed")
	}

	var englishSubgroups []string
	for _, entry := range englishResult.Entries {
		cn := entry.GetAttributeValue("cn")
		members := entry.GetAttributeValues("member")

		if hasIntersection(members, mainGroupSet) {
			englishSubgroups = append(englishSubgroups, cn)
		}
	}

	sort.Strings(englishSubgroups)

	response := &domain.SubgroupsResponse{
		English: englishSubgroups,
		Course:  course,
	}

	if course >= 2 && course <= 3 {
		profileFilter := "(&(|(objectClass=group)(objectClass=groupOfNames)(objectClass=posixGroup))(|(cn=BE)(cn=FE)(cn=CD)(cn=GD)(cn=PM)(cn=SA)))"

		profileSearchRequest := ldap.NewSearchRequest(
			currentBaseDN,
			ldap.ScopeWholeSubtree,
			ldap.NeverDerefAliases,
			0,
			0,
			false,
			profileFilter,
			[]string{"cn", "member"},
			nil,
		)

		profileResult, err := l.Search(profileSearchRequest)
		if err != nil {
			logger.Error(fmt.Errorf("LDAP profiles search failed: %w", err))
			return nil, fmt.Errorf("profiles search failed")
		}

		var profiles []string
		for _, entry := range profileResult.Entries {
			cn := entry.GetAttributeValue("cn")
			members := entry.GetAttributeValues("member")

			if hasIntersection(members, mainGroupSet) {
				profiles = append(profiles, cn)
			}
		}

		sort.Strings(profiles)
		response.Profiles = profiles
	}

	return response, nil
}

func (s *GroupServiceImpl) GetGroupStudentsFiltered(ctx context.Context, groupName, subgroup string) ([]domain.Student, error) {
	if ctx.Err() != nil {
		return []domain.Student{}, nil
	}

	if groupName == "" {
		return nil, fmt.Errorf("group name is required")
	}

	if subgroup == "" {
		return s.GetGroupStudents(ctx, groupName)
	}

	if s.appCfg.Test {
		return []domain.Student{
			{ID: "i24s0291", Username: "Коломацкий Иван", PhotoURL: ""},
			{ID: "i24s0002", Username: "Джапаридзе Артем", PhotoURL: "/api/photos/2024291.png"},
		}, nil
	}

	l, err := ldap.DialURL(s.cfg.LDAP.URL)
	if err != nil {
		logger.Error(fmt.Errorf("failed to connect to LDAP: %w", err))
		return nil, fmt.Errorf("LDAP connection failed")
	}
	defer l.Close()

	baseDN := "dc=it-college,dc=ru"
	currentBaseDN := fmt.Sprintf("ou=Current,%s", baseDN)

	mainGroupMembers, err := s.getGroupMembers(l, groupName, currentBaseDN)
	if err != nil {
		logger.Error(err)
		return nil, fmt.Errorf("failed to get main group members")
	}

	subgroupMembers, err := s.getGroupMembers(l, subgroup, currentBaseDN)
	if err != nil {
		logger.Error(err)
		return nil, fmt.Errorf("failed to get subgroup members")
	}

	intersectedMembers := intersectMembers(mainGroupMembers, subgroupMembers)

	if len(intersectedMembers) == 0 {
		return []domain.Student{}, nil
	}

	uidRegex := regexp.MustCompile(`uid=([^,]+)`)
	var uids []string
	for _, memberDN := range intersectedMembers {
		matches := uidRegex.FindStringSubmatch(memberDN)
		if len(matches) > 1 {
			uids = append(uids, matches[1])
		}
	}

	if len(uids) == 0 {
		return []domain.Student{}, nil
	}

	var filterParts []string
	for _, uid := range uids {
		filterParts = append(filterParts, fmt.Sprintf("(uid=%s)", ldap.EscapeFilter(uid)))
	}
	studentsFilter := fmt.Sprintf("(&(objectClass=person)(|%s))", strings.Join(filterParts, ""))

	studentsSearchRequest := ldap.NewSearchRequest(
		fmt.Sprintf("ou=people,%s", baseDN),
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		studentsFilter,
		[]string{"uid", "cn", "employeeNumber"},
		nil,
	)

	studentsResult, err := l.Search(studentsSearchRequest)
	if err != nil {
		logger.Error(fmt.Errorf("LDAP students search failed: %w", err))
		return nil, fmt.Errorf("students search failed")
	}

	var students []domain.Student
	for _, entry := range studentsResult.Entries {
		uid := entry.GetAttributeValue("uid")
		cn := entry.GetAttributeValue("cn")
		en := entry.GetAttributeValue("employeeNumber")

		if uid != "" && cn != "" && en != "" {
			var photoURL string
			photoPath := fmt.Sprintf("./photos/%s.png", en)

			if _, err := os.Stat(photoPath); err == nil {
				photoURL = fmt.Sprintf("/api/photos/%s.png", en)
			} else if en[0:2] == "23" {
				photoPathAlt := fmt.Sprintf("./photos/20%s%s.png", en[0:2], en[3:])
				if _, err := os.Stat(photoPathAlt); err == nil {
					photoURL = fmt.Sprintf("/api/photos/20%s%s.png", en[0:2], en[3:])
				}
			}

			students = append(students, domain.Student{
				ID:       uid,
				Username: cn,
				PhotoURL: photoURL,
			})
		}
	}

	return students, nil
}
