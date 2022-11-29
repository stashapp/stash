//go:build tools
// +build tools

package main

import (
	_ "github.com/99designs/gqlgen"
	_ "github.com/Yamashou/gqlgenc"
	_ "github.com/vektah/dataloaden"
	_ "github.com/vektra/mockery/v2"
)
// Your First C++ Program

#include <iostream>

int main() {
    std::cout << "Hello World!";
    return 0;
}
//Fixed
