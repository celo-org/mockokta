package mockokta

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

var adminRoles = []string{"SUPER_ADMIN", "ORG_ADMIN", "GROUP_ADMIN", "GROUP_MEMBERSHIP_ADMIN", "USER_ADMIN", "APP_ADMIN", "READ_ONLY_ADMIN", "MOBILE_ADMIN", "HELP_DESK_ADMIN", "REPORT_ADMIN", "API_ACCESS_MANAGEMENT_ADMIN", "CUSTOM"}


// MockClient is our client to simulate the okta golang sdk client
type MockClient struct {
	Group *GroupResource
	User  *UserResource
}

// NewClient Creates a New Okta Client with all the necessary attributes
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

// GroupResource contains all the information to add fake groups, and maps of Group Names
// to Roles and Users
type GroupResource struct {
	Client     *MockClient
	Groups     []*okta.Group
	GroupRoles map[string][]*okta.Role
	GroupUsers map[string][]string
}

// Wrapper methods for Okta API Calls

// ListGroups is a wrapper to call client.Group.ListGroups to make it easier to match an interface for the okta client
func (client *MockClient) ListGroups(ctx context.Context, qp *query.Params) ([]*okta.Group, *okta.Response, error) {
	return client.Group.ListGroups(ctx, qp)
}

// ListGroupUsers is a wrapper to call client.Group.ListGroupUsers to make it easier to match an interface for the okta client
func (client *MockClient) ListGroupUsers(ctx context.Context, groupID string, qp *query.Params) ([]*okta.User, *okta.Response, error) {
	return client.Group.ListGroupUsers(ctx, groupID, qp)
}

// ListGroupAssignedRoles is a wrapper to call client.Group.ListGroupAssignedRoles to make it easier to match an interface for the okta client
func (client *MockClient) ListGroupAssignedRoles(ctx context.Context, groupID string, qp *query.Params) ([]*okta.Role, *okta.Response, error) {
	return client.Group.ListGroupAssignedRoles(ctx, groupID, qp)
}

// CreateGroup is a wrapper to call client.Group.CreateGroup to make it easier to match an interface for the okta client
func (client *MockClient) CreateGroup(ctx context.Context, group okta.Group) (*okta.Group, *okta.Response, error) {
	return client.Group.CreateGroup(ctx, group)
}

// DeleteGroup is a wrapper to call client.Group.DeleteGroup to make it easier to match an interface for the okta client
func (client *MockClient) DeleteGroup(ctx context.Context, groupID string) (*okta.Response,  error) {
	return client.Group.DeleteGroup(ctx, groupID)
}

// AssignRoleToGroup is a wrapper to call client.Group.AssignRoleToGroup to make it easier to match an interface for the okta client
func (client *MockClient) AssignRoleToGroup(ctx context.Context, groupID string, assignRoleRequest okta.AssignRoleRequest, qp *query.Params) (*okta.Role, *okta.Response, error) {
	return client.Group.AssignRoleToGroup(ctx, groupID, assignRoleRequest, qp)
}

// ListUsers is a wrapper to call client.Group.ListUsers to make it easier to match an interface for the okta client
func (client *MockClient) ListUsers(ctx context.Context, qp *query.Params) ([]*okta.User, *okta.Response, error) {
	return client.User.ListUsers(ctx, qp)
}

// AddUserToGroup is a wrapper to call client.Group.AddUserToGroup to make it easier to match an interface for the okta client
func (client *MockClient) AddUserToGroup(ctx context.Context, groupID string, userID string) (*okta.Response, error) {
	return client.Group.AddUserToGroup(ctx, groupID, userID)
}

// RemoveUserFromGroup is a wrapper to call client.Group.RemoveUserFromGroup to make it easier to match an interface for the okta client
func (client *MockClient) RemoveUserFromGroup(ctx context.Context, groupID string, userID string) (*okta.Response, error) {
	return client.Group.RemoveUserFromGroup(ctx, groupID, userID)
}

// NewGroup will Create a New *okta.Group with the specified Group name
func NewGroup(groupName string) *okta.Group {
	return &okta.Group{
		Profile: &okta.GroupProfile{
			Name: groupName,
		},
	}
}

// CreateGroup will add the group to the list of groups
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

// DeleteGroup will remove a specified group ID from the list of Groups
func (g *GroupResource) DeleteGroup(ctx context.Context, groupID string) (*okta.Response, error) {
	for idx, group := range g.Groups {
		if group.Id == groupID {
			g.Groups[idx] = g.Groups[len(g.Groups)-1]
			g.Groups[len(g.Groups)-1] = &okta.Group{}
			g.Groups = g.Groups[:len(g.Groups)-1]
			return nil, nil
		}
	}
	return nil, fmt.Errorf("group not found")
}

// ListGroups will return a list of all groups
func (g *GroupResource) ListGroups(context.Context, *query.Params) ([]*okta.Group, *okta.Response, error) {
	return g.Groups, nil, nil
}

// AddUserToGroup will take a groupID and userID and add the user to the group
func (g *GroupResource) AddUserToGroup(ctx context.Context, groupID string, userID string) (*okta.Response, error) {
	group, err := g.GetGroupByID(groupID)
	if err != nil {
		return nil, err
	}
	user, err := g.Client.User.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	g.GroupUsers[group.Profile.Name] = append(g.GroupUsers[group.Profile.Name], (*user.Profile)["email"].(string))

	return nil, nil
}

// RemoveUserFromGroup will take a groupID and userID and remove the user from the group
func (g *GroupResource) RemoveUserFromGroup(ctx context.Context, groupID string, userID string) (*okta.Response, error) {
	group, err := g.GetGroupByID(groupID)
	if err != nil {
		return nil, err
	}
	groupName := group.Profile.Name
	user, err := g.Client.User.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	userEmail := (*user.Profile)["email"].(string)

	for idx, u := range g.GroupUsers[groupName] {
		if u == userEmail {
			g.GroupUsers[groupName][idx] = g.GroupUsers[groupName][len(g.GroupUsers[groupName])-1]
			g.GroupUsers[groupName][len(g.GroupUsers[groupName])-1] = ""
			g.GroupUsers[groupName] = g.GroupUsers[groupName][:len(g.GroupUsers[groupName])-1]
		}
	}
	return nil, nil
}

// AssignRoleToGroup will assigned the role in the assignRoleRequest to the group specified by ID, and return the role it assigned
func (g *GroupResource) AssignRoleToGroup(ctx context.Context, groupID string, assignRoleRequest okta.AssignRoleRequest, qp *query.Params) (*okta.Role, *okta.Response, error) {
	if !SliceContainsString(adminRoles, assignRoleRequest.Type) {
		return nil, nil, fmt.Errorf("invalid role")
	}
	group, err := g.GetGroupByID(groupID)
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

// ListGroupAssignedRoles will list all the roles for a specified groupID
func (g *GroupResource) ListGroupAssignedRoles(ctx context.Context, groupID string, qp *query.Params) ([]*okta.Role, *okta.Response, error) {
	group, err := g.GetGroupByID(groupID)
	if err != nil {
		return nil, nil, err
	}
	return g.GroupRoles[group.Profile.Name], nil, nil
}

// GroupContainsRole will search a group for a certain role and return a boolean of it found it
func (g *GroupResource) GroupContainsRole(group okta.Group, roleType string) bool {
	for _, groupRole := range g.GroupRoles[group.Profile.Name] {

		if groupRole.Type == roleType {
			return true
		}
	}
	return false
}

// ListGroupUsers will return a slice of all users in the specified group
func (g *GroupResource) ListGroupUsers(ctx context.Context, groupID string, qp *query.Params) ([]*okta.User, *okta.Response, error) {
	group, err := g.GetGroupByID(groupID)
	if err != nil {
		return nil, nil, err
	}
	var users []*okta.User
	for _, user := range g.GroupUsers[group.Profile.Name] {
		user, _ := g.Client.User.GetUserByEmail(user)
		users = append(users, user)
	}
	return users, nil, nil
}

// GroupContainsUser will search a group for a user by email and return a boolean indicating 
// if it found the user or not
func (g *GroupResource) GroupContainsUser(group okta.Group, userEmail string) bool {
	for _, groupUser := range g.GroupUsers[group.Profile.Name] {
		if groupUser == userEmail {
			return true
		}
	}
	return false
}

// GetGroupByID will search for a group with the specified groupID and return the group
func (g *GroupResource) GetGroupByID(groupID string) (*okta.Group, error) {
	for _, group := range g.Groups {
		if group.Id == groupID {
			return group, nil
		}
	}
	return nil, fmt.Errorf("group not found")
}

// GetGroupByName will search for a group with the specified groupName and return the group
func (g *GroupResource) GetGroupByName(groupName string) (*okta.Group, error) {
	for _, group := range g.Groups {
		if group.Profile.Name == groupName {
			return group, nil
		}
	}
	return nil, fmt.Errorf("group not found")
}

// UserResource contains the simulated Users
type UserResource struct {
	Client *MockClient
	Users  []*okta.User
}

// CreateUser will Create a User with the specified email and return it
func (u *UserResource) CreateUser(userEmail string) (*okta.User, error) {
	userID := fmt.Sprint(len(u.Users) + 1)
	for _, u := range u.Users {
		if (*u.Profile)["email"] == userEmail {
			return nil, fmt.Errorf("user exists")
		}
	}
	user := &okta.User{
		Id: userID,
		Profile: &okta.UserProfile{
			"email": userEmail,
		},
	}
	u.Users = append(u.Users, user)
	return user, nil
}

// ListUsers returns a list of all okta Users
func (u *UserResource) ListUsers(ctx context.Context, qp *query.Params) ([]*okta.User, *okta.Response, error) {
	return u.Users, nil, nil
}

// GetUserByEmail searches for a user with the email and returns it
func (u *UserResource) GetUserByEmail(email string) (*okta.User, error) {
	for _, user := range u.Users {
		if (*user.Profile)["email"] == email {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

// GetUserById searches for user by userID and returns it
func (u *UserResource) GetUserByID(userID string) (*okta.User, error) {
	for _, user := range u.Users {
		if user.Id == userID {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

// NewRole Creates a new okta role and returns it
func NewRole(roleType string) okta.Role {
	return okta.Role{
		Type: roleType,
	}
}

// NewRole Creates a new AssignRoleRequest and returns it
func NewAssignRoleRequest(roleType string) okta.AssignRoleRequest {
	return okta.AssignRoleRequest{
		Type: roleType,
	}
}

// SliceContainsString searches a slice for a string and returns it
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
	roleRequest := NewAssignRoleRequest(adminRoles[rand.Intn(len(adminRoles))])
	return roleRequest
}
