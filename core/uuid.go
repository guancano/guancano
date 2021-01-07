package core

import "github.com/rogpeppe/fastuuid"

// global generator for UUIDs
var generator, _ = fastuuid.NewGenerator()
