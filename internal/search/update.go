package search

import (
	"context"
	"path"

	"github.com/alist-org/alist/v3/internal/model"
	"github.com/alist-org/alist/v3/internal/op"
	mapset "github.com/deckarep/golang-set/v2"
	log "github.com/sirupsen/logrus"
)

func Update(parent string, objs []model.Obj) {
	if !instance.Config().AutoUpdate {
		return
	}
	ctx := context.Background()
	nodes, err := instance.Get(ctx, parent)
	if err != nil {
		log.Errorf("update search index error while get nodes: %+v", err)
		return
	}
	now := mapset.NewSet[string]()
	for i := range objs {
		now.Add(objs[i].GetName())
	}
	old := mapset.NewSet[string]()
	for i := range nodes {
		old.Add(nodes[i].Name)
	}
	// delete data that no longer exists
	toDelete := old.Difference(now)
	toAdd := now.Difference(old)
	for i := range nodes {
		if toDelete.Contains(nodes[i].Name) {
			err = instance.Del(ctx, path.Join(parent, nodes[i].Name))
			if err != nil {
				log.Errorf("update search index error while del old node: %+v", err)
				return
			}
		}
	}
	for i := range objs {
		if toAdd.Contains(objs[i].GetName()) {
			err = instance.Index(ctx, parent, objs[i])
			if err != nil {
				log.Errorf("update search index error while index new node: %+v", err)
				return
			}
		}
	}
}

func init() {
	op.RegisterObjsUpdateHook(Update)
}