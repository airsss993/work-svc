package service

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/airsss993/work-svc/internal/config"
	"github.com/airsss993/work-svc/internal/domain"
	"github.com/airsss993/work-svc/pkg/logger"
	"github.com/go-ldap/ldap/v3"
)

type GroupService interface {
	GetITGroups(ctx context.Context) ([]domain.GroupInfo, error)
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
