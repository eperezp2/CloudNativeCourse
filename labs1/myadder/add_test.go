package myadder

import "testing"

fun TestAdd(t *testing.T){
	want := 7
	got := Add(3,4)
	if want !=got{
		t.Errorf("Error in myadder; Want 7, got %d", got)
		
	}
}