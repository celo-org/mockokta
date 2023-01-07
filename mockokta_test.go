package mockokta

import (
	"context"
	"math/rand"
	"reflect"
	"testing"
	"time"

//	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/stretchr/testify/assert"
)

func TestGroupResource_CreateGroup(t *testing.T) {
	t.Run("should not create group with empty name", func(t *testing.T) {
		groupNameArg := ""
		client := NewClient()
		group := NewGroup(groupNameArg)

		_, _, err := client.Group.CreateGroup(context.TODO(), *group)

		if err == nil {
			t.Errorf("expected error but didn't get one")
		}
	})

	t.Run("should not create group with invalid length", func(t *testing.T) {
		groupNameArg := RandStringRunes(256)
		client := NewClient()
		group := NewGroup(groupNameArg)

		_, _, err := client.Group.CreateGroup(context.TODO(), *group)

		if err == nil {
			t.Errorf("expected error but didn't get one")
		}
	})
	t.Run("should not create group if it already exists", func(t *testing.T) {
		groupNameArg := "TestGroup"
		client := NewClient()
		group := NewGroup(groupNameArg)

		client.Group.CreateGroup(context.TODO(), *group)
		_, _, err := client.Group.CreateGroup(context.TODO(), *group)

		if err == nil {
			t.Errorf("expected error but didn't get one")
		}
	})

	t.Run("should create group", func(t *testing.T) {
		groupNameArg := "TestGroup"
		client := NewClient()
		group := NewGroup(groupNameArg)

		client.Group.CreateGroup(context.TODO(), *group)

		want := group
		got := client.Group.Groups[0]

		if reflect.DeepEqual(got, want) {
			t.Fatalf("got %v want %v", got, want)
		}
	})
}

func TestGroupResource_ListGroups(t *testing.T) {
	group1NameArg := "TestGroup1"
	group2NameArg := "TestGroup2"

	client := NewClient()
	group1 := NewGroup(group1NameArg)
	group2 := NewGroup(group2NameArg)

	client.Group.Groups = append(client.Group.Groups, group1, group2)

	want := []*okta.Group{group1, group2}
	got, _, _ := client.Group.ListGroups(context.TODO(), nil)

	assert.ElementsMatch(t, got, want)
}

func TestGroupResource_AssignRoleToGroup(t *testing.T) {
	t.Run("should not assign role with invalid name to group", func(t *testing.T) {
		groupNameArg := "TestGroup"
		roleArg := "Invalid_Role"
        roleRequest := NewAssignRoleRequest(roleArg)

		client := NewClient()
		group, _, _ := client.Group.CreateGroup(context.TODO(), *NewGroup(groupNameArg))

		_, _, err := client.Group.AssignRoleToGroup(context.TODO(), group.Id, roleRequest, nil)

		if err == nil {
			t.Errorf("Expected err, but didn't get one")
		}
	})

	t.Run("should return error if group doesn't exist", func(t *testing.T) {
        groupNameArg := "NonexistentGroup"
       
		client := NewClient()
		roleRequest := RandAdminRoleRequest()

		_, _, err := client.Group.AssignRoleToGroup(context.TODO(), groupNameArg, roleRequest, nil)

		if err == nil {
			t.Errorf("Expected err, but didn't get one")
		}
	})
	t.Run("should return error if role exists for group", func(t *testing.T) {
		groupNameArg := "TestGroup"
		roleRequest := RandAdminRoleRequest()

		client := NewClient()
		group, _, _ := client.Group.CreateGroup(context.TODO(), *NewGroup(groupNameArg))

		client.Group.AssignRoleToGroup(context.TODO(), group.Id, roleRequest, nil)
		_, _, err := client.Group.AssignRoleToGroup(context.TODO(), group.Id, roleRequest, nil)

		if err == nil {
			t.Errorf("Expected err, but didn't get one")
		}
	})

	t.Run("should assign role to group", func(t *testing.T) {
		groupNameArg := "TestGroup"
		roleRequest := RandAdminRoleRequest()
		client := NewClient()
		group, _, _ := client.Group.CreateGroup(context.TODO(), *NewGroup(groupNameArg))
		client.Group.AssignRoleToGroup(context.TODO(), group.Id, roleRequest, nil)

        got := client.Group.GroupContainsRole(*group, roleRequest.Type)

        if got != true {
            t.Errorf("expected group %+v to contain role %v but was not found", client.Group, roleRequest.Type)
        }
	})
}

func TestGroupResource_ListGroupAssignedRoles(t *testing.T) {
	t.Run("should return error if the group doesn't exist", func(t *testing.T) {
		groupIdArg := "NonexistentGroup"
		roleRequest := RandAdminRoleRequest()

		client := NewClient()
	    _, _, err := client.Group.AssignRoleToGroup(context.TODO(), groupIdArg, roleRequest, nil)
       
		if err == nil {
			t.Errorf("Expected err, but didn't get one")
		}

	})

	t.Run("should list assigned roles", func(t *testing.T) {
		groupNameArg := "TestGroup"

		roleTypeArg1 := "SUPER_ADMIN"
        roleTypeArg2 := "GROUP_ADMIN"

        roleRequest1 := NewAssignRoleRequest(roleTypeArg1)
        roleRequest2 := NewAssignRoleRequest(roleTypeArg2)
       
		client := NewClient()
		group, _, _ := client.Group.CreateGroup(context.TODO(), *NewGroup(groupNameArg))

		role1, _, _ :=  client.Group.AssignRoleToGroup(context.TODO(), group.Id, roleRequest1, nil)
		role2, _, _ := client.Group.AssignRoleToGroup(context.TODO(), group.Id, roleRequest2, nil)
       
        want := []*okta.Role{role1, role2}
        got, _, _ := client.Group.ListGroupAssignedRoles(context.TODO(), group.Id, nil)

        if !reflect.DeepEqual(got, want) {
            t.Errorf("expected roles %v to match %v", got, want)
        }
	})
}

func TestUserResource_CreateUser(t *testing.T) {
    t.Run("should err if user exists", func(t *testing.T) {
        userEmail := "TestUser@test.com"

        client := NewClient()

        client.User.CreateUser(userEmail)
        _, err := client.User.CreateUser(userEmail)

 		if err == nil {
			t.Errorf("expected error but didn't get one")
		}
    })

    t.Run("should create user", func(t *testing.T) {
        userEmail := "TestUser@test.com"

        client := NewClient()

        want, _ := client.User.CreateUser(userEmail)
        got, _ := client.User.GetUserByEmail(context.TODO(), userEmail)
        
        if !reflect.DeepEqual(got, want) {
            t.Errorf("got %v want %v", got, want)
        }
    })
}

func TestUserResource_ListUsers(t *testing.T) {
    client := NewClient()
    userEmailArg1 := "TestUser1"
    userEmailArg2 := "TestUser2"

    user1, _ := client.User.CreateUser(userEmailArg1)
    user2, _ := client.User.CreateUser(userEmailArg2)

    want := []*okta.User{user1, user2}
    got, _, _ := client.User.ListUsers(context.TODO(), nil)

    
    if !reflect.DeepEqual(got, want) {
        t.Errorf("got %v want %v", got, want)
    }
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}


