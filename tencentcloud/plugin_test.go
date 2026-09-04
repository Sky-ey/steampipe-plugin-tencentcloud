package tencentcloud

import (
	"context"
	"sort"
	"testing"

	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
)

func TestTableColumnNamesUnique(t *testing.T) {
	for tableName, table := range Plugin(context.Background()).TableMap {
		seen := make(map[string]int, len(table.Columns))
		for _, column := range table.Columns {
			if column == nil {
				t.Errorf("表 %s 存在 nil 列定义", tableName)
				continue
			}
			seen[column.Name]++
		}

		duplicates := make([]string, 0)
		for name, count := range seen {
			if count > 1 {
				duplicates = append(duplicates, name)
			}
		}
		sort.Strings(duplicates)

		if len(duplicates) > 0 {
			t.Errorf("表 %s 存在重复列 %v", tableName, duplicates)
		}
	}
}

func TestCommonColumnsPresentOnEveryTable(t *testing.T) {
	commonNames := make([]string, 0)
	for _, column := range utils.CommonColumns() {
		commonNames = append(commonNames, column.Name)
	}

	for tableName, table := range Plugin(context.Background()).TableMap {
		present := make(map[string]bool, len(table.Columns))
		for _, column := range table.Columns {
			if column != nil {
				present[column.Name] = true
			}
		}
		for _, name := range commonNames {
			if !present[name] {
				t.Errorf("表 %s 缺少公共列 %s——请使用 utils.WithCommonColumns 构建列定义", tableName, name)
			}
		}
	}
}

func TestConnectionKeyColumnsAreDeclared(t *testing.T) {
	p := Plugin(context.Background())

	for _, keyColumn := range p.ConnectionKeyColumns {
		if keyColumn.Hydrate == nil {
			t.Errorf("连接键列 %s 未配置 Hydrate 函数", keyColumn.Name)
		}
		for tableName, table := range p.TableMap {
			found := false
			for _, column := range table.Columns {
				if column != nil && column.Name == keyColumn.Name {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("表 %s 缺少连接键列 %s", tableName, keyColumn.Name)
			}
		}
	}
}
