// Package iter provides logic for one dimensional iterators.
//
// Anytime a channel iterator is used it is important to use all the values or
// it will leak go routines.
package iter

// BufferLen is the buffer length used for channels
var BufferLen = 100
