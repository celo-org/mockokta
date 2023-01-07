package mockokta

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

var ADMIN_ROLES = []string{"SUPER_ADMIN", "ORG_ADMIN", "GROUP_ADMIN", "GROUP_MEMBERSHIP_ADMIN", "USER_ADMIN", "APP_ADMIN", "READ_ONLY_ADMIN", "MOBILE_ADMIN", "HELP_DESK_ADMIN", "REPORT_ADMIN", "API_ACCESS_MANAGEMENT_ADMIN", "CUSTOM"}

// func NewClient(ctx context.Context, conf ...ConfigSetter) (context.Context, *Client, error) {
type MockClient struct {
	Group *GroupResource
	User  *UserResource
}

func NewClient() *MockClient {
	c := &MockClient{}
	c.Group = &GroupResource{
		Client:     c,
		GroupRoles: make(map[string][]*okta.Role),
        GroupUsers: make(map[string][]string),
	}
	c.User = &UserResource{
		Client: c,
	}
	return c
}

type GroupResource struct {
	Client     *MockClient
	Groups     []*okta.Group
	GroupRoles map[string][]*okta.Role
	GroupUsers map[string][]string
}

// Wrapper methods for Okta API Calls
func (client *MockClient) ListGroups(ctx context.Context, qp *query.Params) ([]*okta.Group, *okta.Response, error) {
	return client.Group.ListGroups(ctx, qp)
}

func (client *MockClient) ListGroupAssignedRoles(ctx context.Context, groupId string, qp *query.Params) ([]*okta.Role, *okta.Response, error) {
	return client.Group.ListGroupAssignedRoles(ctx, groupId, qp)
}

func (client *MockClient) CreateGroup(ctx context.Context, group okta.Group) (*okta.Group, *okta.Response, error) {
	return client.Group.CreateGroup(ctx, group)
}

func (client *MockClient) AssignRoleToGroup(ctx context.Context, groupId string, assignRoleRequest okta.AssignRoleRequest, qp *query.Params) (*okta.Role, *okta.Response, error) {
	return client.Group.AssignRoleToGroup(ctx, groupId, assignRoleRequest, qp)
}

func (client *MockClient) ListUsers(ctx context.Context, qp *query.Params) ([]*okta.User, *okta.Response, error) {
	return client.User.ListUsers(ctx, qp)
}

func (client *MockClient) AddUserToGroup(ctx context.Context, groupId string, userId string) (*okta.Response, error) {
	return client.Group.AddUserToGroup(ctx, groupId, userId)
}

func (g *GroupResource) CreateGroup(ctx context.Context, group okta.Group) (*okta.Group, *okta.Response, error) {

	group.Id = fmt.Sprint(len(g.Groups) + 1)
	for _, x := range g.Groups {

		if x.Profile.Name == group.Profile.Name {
			return nil, nil, fmt.Errorf("unable to create group: group exists")
		}
	}

	if len(group.Profile.Name) > 255 || len(group.Profile.Name) < 1 {
		return nil, nil, fmt.Errorf("unable to create group: invalid name length")
	}

	g.Groups = append(g.Groups, &group)
	return &group, nil, nil
}

func (g *GroupResource) ListGroups(context.Context, *query.Params) ([]*okta.Group, *okta.Response, error) {
	return g.Groups, nil, nil
}

func NewGroup(groupName string) *okta.Group {
	return &okta.Group{
		Profile: &okta.GroupProfile{
			Name: groupName,
		},
	}
}

func (g *GroupResource) AddUserToGroup(ctx context.Context, groupId string, userId string) (*okta.Response, error) {
	group, err := g.GetGroupById(groupId)
	if err != nil {
		return nil, err
	}
	user, err := g.Client.User.GetUserById(userId)
	if err != nil {
		return nil, err
	}
	g.GroupUsers[group.Profile.Name] = append(g.GroupUsers[group.Profile.Name], (*user.Profile)["email"].(string))
	return nil, nil
}

func (g *GroupResource) AssignRoleToGroup(ctx context.Context, groupId string, assignRoleRequest okta.AssignRoleRequest, qp *query.Params) (*okta.Role, *okta.Response, error) {
	if !SliceContainsString(ADMIN_ROLES, assignRoleRequest.Type) {
		return nil, nil, fmt.Errorf("invalid role")
	}
	group, err := g.GetGroupById(groupId)
	if err != nil {
		return nil, nil, err
	}

	if g.GroupContainsRole(*group, assignRoleRequest.Type) {
		return nil, nil, fmt.Errorf("group role exists")
	}

	role := NewRole(assignRoleRequest.Type)
	role.Id = fmt.Sprintf("%v", len(g.GroupRoles)+1)
	g.GroupRoles[group.Profile.Name] = append(g.GroupRoles[group.Profile.Name], &role)
	return &role, nil, nil
}

func (g *GroupResource) ListGroupAssignedRoles(ctx context.Context, groupId string, qp *query.Params) ([]*okta.Role, *okta.Response, error) {
	group, err := g.GetGroupById(groupId)
	if err != nil {
		return nil, nil, err
	}

	return g.GroupRoles[group.Profile.Name], nil, nil
}

func (g *GroupResource) GroupContainsRole(group okta.Group, roleType string) bool {
	for _, groupRole := range g.GroupRoles[group.Profile.Name] {

		if groupRole.Type == roleType {
			return true
		}
	}
	return false
}

func (g *GroupResource) GroupContainsUser(group okta.Group, userEmail string) bool {
	for _, groupUser := range g.GroupUsers[group.Profile.Name] {
        if groupUser == userEmail {
			return true
		}
	}
	return false
}

func (g *GroupResource) GetGroupById(groupId string) (*okta.Group, error) {
	for _, group := range g.Groups {
		if group.Id == groupId {
			return group, nil
		}
	}
	return nil, fmt.Errorf("group not found")
}

type UserResource struct {
	Client *MockClient
	Users  []*okta.User
}

func (u *UserResource) CreateUser(userEmail string) (*okta.User, error) {
	userId := fmt.Sprint(len(u.Users) + 1)

	for _, u := range u.Users {
		if (*u.Profile)["email"] == userEmail {
			return nil, fmt.Errorf("user exists")
		}
	}

	user := &okta.User{
		Id: userId,
		Profile: &okta.UserProfile{
			"email": userEmail,
		},
	}

	u.Users = append(u.Users, user)
	return user, nil
}

func (u *UserResource) ListUsers(ctx context.Context, qp *query.Params) ([]*okta.User, *okta.Response, error) {
	return u.Users, nil, nil
}

func (u *UserResource) GetUserByEmail(email string) (*okta.User, error) {
	for _, user := range u.Users {
		if (*user.Profile)["email"] == email {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (u *UserResource) GetUserById(userId string) (*okta.User, error) {
	for _, user := range u.Users {
		if user.Id == userId {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func NewRole(roleType string) okta.Role {
	return okta.Role{
		Type: roleType,
	}
}

func NewAssignRoleRequest(roleType string) okta.AssignRoleRequest {
	return okta.AssignRoleRequest{
		Type: roleType,
	}
}

func SliceContainsString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

func RandAdminRoleRequest() okta.AssignRoleRequest {
	rand.Seed(time.Now().UnixNano())
	roleRequest := NewAssignRoleRequest(ADMIN_ROLES[rand.Intn(len(ADMIN_ROLES))])
	return roleRequest
}
