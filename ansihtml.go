// Package ansihtml parses text formatted with ANSI escape sequences and
// outputs text suitable for display in a HTML <pre> tag.
//
// Text effects are encoded as <span> tags with various classes:
//
// - Foreground colors in 3- or 8-bit codes become classes 'fg0' .. 'fg255'.
// - Background colors in 3- or 8-bit codes become classes 'bg0' .. 'bg255'.
// - Bold becomes class 'bold'.
// - Faint becomes class 'faint'.
package ansihtml
