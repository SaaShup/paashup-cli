package main

import (
    "testing"
)

func TestListContainers(t *testing.T) {
    _, err := listContainers()
    if err != nil {
        t.Errorf("Error: %s", err)
    }
}

func TestListContainersByHost(t *testing.T) {
    _, err := listContainersByHost(1)
    if err != nil {
        t.Errorf("Error: %s", err)
    }
}

