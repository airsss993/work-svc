package service

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/airsss993/work-svc/internal/config"
	"github.com/airsss993/work-svc/internal/domain"
	"github.com/airsss993/work-svc/pkg/logger"
	"github.com/go-ldap/ldap/v3"
)

type GroupService interface {
	GetITGroups(ctx context.Context) ([]domain.GroupInfo, error)
	GetGroupStudents(ctx context.Context, groupName string) ([]domain.Student, error)
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
		fmt.Sprintf("ou=groups,%s", baseDN),
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
		fmt.Sprintf("ou=groups,%s", baseDN),
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
			photoURL := ""
			photoPath := fmt.Sprintf("./photos/%s.png", en)

			if _, err := os.Stat(photoPath); err == nil {
				photoURL = fmt.Sprintf("/api/photos/%s.png", en)
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
