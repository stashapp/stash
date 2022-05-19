package singlefile

import (
	"context"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/require"
)

type UpdatePtrToPtrResults struct {
	UpdatedPtrToPtr PtrToPtrOuter `json:"updatePtrToPtr"`
}

func TestPtrToPtr(t *testing.T) {
	resolvers := &Stub{}

	c := client.New(handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: resolvers})))

	resolvers.MutationResolver.UpdatePtrToPtr = func(ctx context.Context, in UpdatePtrToPtrOuter) (ret *PtrToPtrOuter, err error) {
		ret = &PtrToPtrOuter{
			Name: "oldName",
			Inner: &PtrToPtrInner{
				Key:   "oldKey",
				Value: "oldValue",
			},
			StupidInner: nest7(&PtrToPtrInner{
				Key:   "oldStupidKey",
				Value: "oldStupidValue",
			}),
		}

		if in.Name != nil {
			ret.Name = *in.Name
		}

		if in.Inner != nil {
			inner := *in.Inner
			if inner == nil {
				ret.Inner = nil
			} else {
				if in.Inner == nil {
					ret.Inner = &PtrToPtrInner{}
				}
				if inner.Key != nil {
					ret.Inner.Key = *inner.Key
				}
				if inner.Value != nil {
					ret.Inner.Value = *inner.Value
				}
			}
		}

		if in.StupidInner != nil {
			si := *in.StupidInner
			if si == nil {
				ret.StupidInner = nil
			} else {
				deepIn := ******si
				deepOut := ******ret.StupidInner
				if deepIn.Key != nil {
					deepOut.Key = *deepIn.Key
				}
				if deepIn.Value != nil {
					deepOut.Value = *deepIn.Value
				}
			}
		}
		return
	}

	t.Run("pointer to pointer input missing", func(t *testing.T) {
		var resp UpdatePtrToPtrResults

		err := c.Post(`mutation { updatePtrToPtr(input: { name: "newName" }) { name, inner { key, value }, stupidInner { key, value }}}`, &resp)
		require.NoError(t, err)

		require.Equal(t, resp.UpdatedPtrToPtr.Name, "newName")
		require.NotNil(t, resp.UpdatedPtrToPtr.Inner)
		require.Equal(t, resp.UpdatedPtrToPtr.Inner.Key, "oldKey")
		require.Equal(t, resp.UpdatedPtrToPtr.Inner.Value, "oldValue")
		require.NotNil(t, resp.UpdatedPtrToPtr.StupidInner)
		require.NotNil(t, ******resp.UpdatedPtrToPtr.StupidInner)
		require.Equal(t, (******resp.UpdatedPtrToPtr.StupidInner).Key, "oldStupidKey")
		require.Equal(t, (******resp.UpdatedPtrToPtr.StupidInner).Value, "oldStupidValue")
	})

	t.Run("pointer to pointer input non-null", func(t *testing.T) {
		var resp UpdatePtrToPtrResults

		err := c.Post(`mutation {
			updatePtrToPtr(input: {
				inner: {
					key: "newKey"
					value: "newValue"
				}
			})
			{ name, inner { key, value }, stupidInner { key, value }}
		}`, &resp)
		require.NoError(t, err)

		require.Equal(t, resp.UpdatedPtrToPtr.Name, "oldName")
		require.NotNil(t, resp.UpdatedPtrToPtr.Inner)
		require.Equal(t, resp.UpdatedPtrToPtr.Inner.Key, "newKey")
		require.Equal(t, resp.UpdatedPtrToPtr.Inner.Value, "newValue")
		require.NotNil(t, resp.UpdatedPtrToPtr.StupidInner)
		require.NotNil(t, ******resp.UpdatedPtrToPtr.StupidInner)
		require.Equal(t, (******resp.UpdatedPtrToPtr.StupidInner).Key, "oldStupidKey")
		require.Equal(t, (******resp.UpdatedPtrToPtr.StupidInner).Value, "oldStupidValue")
	})

	t.Run("pointer to pointer input null", func(t *testing.T) {
		var resp UpdatePtrToPtrResults

		err := c.Post(`mutation { updatePtrToPtr(input: { inner: null }) { name, inner { key, value }, stupidInner { key, value }}}`, &resp)
		require.NoError(t, err)

		require.Equal(t, resp.UpdatedPtrToPtr.Name, "oldName")
		require.Nil(t, resp.UpdatedPtrToPtr.Inner)
		require.NotNil(t, resp.UpdatedPtrToPtr.StupidInner)
		require.NotNil(t, ******resp.UpdatedPtrToPtr.StupidInner)
		require.Equal(t, (******resp.UpdatedPtrToPtr.StupidInner).Key, "oldStupidKey")
		require.Equal(t, (******resp.UpdatedPtrToPtr.StupidInner).Value, "oldStupidValue")
	})

	t.Run("many pointers input non-null", func(t *testing.T) {
		var resp UpdatePtrToPtrResults

		err := c.Post(`mutation {
			updatePtrToPtr(input: {
				stupidInner: {
					key: "newKey"
					value: "newValue"
				}
			})
			{ name, inner { key, value }, stupidInner { key, value }}
		}`, &resp)
		require.NoError(t, err)

		require.Equal(t, resp.UpdatedPtrToPtr.Name, "oldName")
		require.NotNil(t, resp.UpdatedPtrToPtr.Inner)
		require.Equal(t, resp.UpdatedPtrToPtr.Inner.Key, "oldKey")
		require.Equal(t, resp.UpdatedPtrToPtr.Inner.Value, "oldValue")
		require.NotNil(t, resp.UpdatedPtrToPtr.StupidInner)
		require.NotNil(t, ******resp.UpdatedPtrToPtr.StupidInner)
		require.Equal(t, (******resp.UpdatedPtrToPtr.StupidInner).Key, "newKey")
		require.Equal(t, (******resp.UpdatedPtrToPtr.StupidInner).Value, "newValue")
	})

	t.Run("many pointers input null", func(t *testing.T) {
		var resp UpdatePtrToPtrResults

		err := c.Post(`mutation { updatePtrToPtr(input: { stupidInner: null }) { name, inner { key, value }, stupidInner { key, value }}}`, &resp)
		require.NoError(t, err)

		require.Equal(t, resp.UpdatedPtrToPtr.Name, "oldName")
		require.NotNil(t, resp.UpdatedPtrToPtr.Inner)
		require.Equal(t, resp.UpdatedPtrToPtr.Inner.Key, "oldKey")
		require.Equal(t, resp.UpdatedPtrToPtr.Inner.Value, "oldValue")
		require.Nil(t, resp.UpdatedPtrToPtr.StupidInner)
	})
}

func nest7(in *PtrToPtrInner) *******PtrToPtrInner {
	si2 := &in
	si3 := &si2
	si4 := &si3
	si5 := &si4
	si6 := &si5
	si7 := &si6

	return si7
}
