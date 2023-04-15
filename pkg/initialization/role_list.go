package initialization

import (
	"errors"
	"log"
	"os"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/duke-git/lancet/v2/validator"
	"gopkg.in/yaml.v2"
)

// Export Global

var roleList RoleList

func GetAllUniqueTags() *[]string {
	return roleList.GetAllUniqueTags()
}

func GetRoleByTitle(title string) *Role {
	return roleList.GetRoleByTitle(title)
}

func GetTitleListByTag(tags string) *[]string {
	return roleList.GetTitleListByTag(tags)
}

func GetFirstRoleContentByTitle(title string) (string, error) {
	return roleList.GetFirstRoleContentByTitle(title)
}

// InitRoleList 加载Prompt
func InitRoleList() RoleList {
	data, err := os.ReadFile("role_list.yaml")
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(data, &roleList)
	if err != nil {
		log.Fatal(err)
	}
	return roleList
}

type Role struct {
	Title   string   `yaml:"title"`
	Content string   `yaml:"content"`
	Tags    []string `yaml:"tags"`
}

type RoleList []Role

func (rl RoleList) GetAllUniqueTags() *[]string {
	tags := make([]string, 0)
	for _, role := range rl {
		tags = append(tags, role.Tags...)
	}
	result := slice.Union(tags)
	return &result
}

func (rl RoleList) GetRoleByTitle(title string) *Role {
	for _, role := range rl {
		if role.Title == title {
			return &role
		}
	}
	return nil
}

func (rl RoleList) GetTitleListByTag(tags string) *[]string {
	roles := make([]string, 0)
	//pp.Println(RoleList)
	for _, role := range rl {
		for _, roleTag := range role.Tags {
			if roleTag == tags && !validator.IsEmptyString(role.
				Title) {
				roles = append(roles, role.Title)
			}
		}
	}
	return &roles
}

func (rl RoleList) GetFirstRoleContentByTitle(title string) (string, error) {
	for _, role := range rl {
		if role.Title == title {
			return role.Content, nil
		}
	}
	return "", errors.New("role not found")
}
