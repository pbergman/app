package app

import (
    "sort"
)

const (
    DEFAULT_GROUP_NAME = "000000"
)

type ByName []CommandInterface

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].GetName() < a[j].GetName() }

type commandList struct {
    Max      int
    Runnable map[string][]CommandInterface
    Helpers  map[string][]CommandInterface
}

func (c *commandList) NotDefaultGroup(s string) bool {
    return s != DEFAULT_GROUP_NAME
}

func (c *commandList) HasCommands(runnable bool) bool {
    if runnable {
        return len(c.Runnable) > 0
    } else {
        return len(c.Helpers) > 0
    }
}

func (c *commandList) GetGroups(runnable bool) []string {
    var list []string
    if runnable {
        for key := range c.Runnable {
            list = append(list, key)
        }
    } else {
        for key := range c.Helpers {
            list = append(list, key)
        }
    }
    sort.Strings(list)
    return list
}

func (c *commandList) GetCommands(runnable bool, group string) []CommandInterface {
    var list []CommandInterface
    if runnable {
        list = c.Runnable[group]
    } else {
        list = c.Helpers[group]
    }
    sort.Sort(ByName(list))
    return list
}