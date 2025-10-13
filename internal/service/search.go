package service

import (
	"context"
	"fmt"
	"os"

	"github.com/airsss993/work-svc/internal/config"
	"github.com/airsss993/work-svc/internal/domain"
	"github.com/airsss993/work-svc/pkg/logger"
	"github.com/go-ldap/ldap/v3"
)

type StudentService interface {
	SearchStudents(ctx context.Context, query string) ([]domain.Student, error)
}

type StudentServiceImpl struct {
	cfg    *config.Config
	appCfg *config.App
}

func NewStudentService(cfg *config.Config, appCfg *config.App) *StudentLDAPService {
	return &StudentLDAPService{
		cfg:    cfg,
		appCfg: appCfg,
	}
}

func (s *StudentServiceImpl) SearchStudents(ctx context.Context, query string) ([]domain.Student, error) {
	if ctx.Err() != nil {
		return []domain.Student{}, nil
	}

	if query == "" {
		return []domain.Student{}, nil
	}

	if s.appCfg.Test {
		return []domain.Student{
			{ID: "i24s0291", Username: "Коломацкий Иван"},
			{ID: "i24s0002", Username: "Джапаридзе Артем"},
		}, nil
	}

	l, err := ldap.DialURL(s.cfg.LDAP.URL)
	if err != nil {
		logger.Error(fmt.Errorf("failed to connect to LDAP: %w", err))
		return nil, fmt.Errorf("LDAP connection failed")
	}
	defer l.Close()

	filter := fmt.Sprintf(
		"(&(objectClass=person)(!(uid=t*))(|(uid=*%s*)(cn=*%s*)))",
		ldap.EscapeFilter(query),
		ldap.EscapeFilter(query),
	)

	searchRequest := ldap.NewSearchRequest(
		"ou=people,dc=it-college,dc=ru",
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		50,
		0,
		false,
		filter,
		[]string{"uid", "cn", "employeeNumber"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		logger.Error(fmt.Errorf("LDAP search failed: %w", err))
		return nil, fmt.Errorf("search failed")
	}

	var students []domain.Student
	for _, entry := range sr.Entries {
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
