package notes

import (
	"fmt"
)

// Note represents a single note that can be interacted with by the user.
type Note struct {
	NoteID  string
	Message string
}

// NotePage - a page of notes.
type NotePage struct {
	Notes []Note
}

var NoteList = make([]Note, 0)

func getNote(userID string, noteID string) (Note, error) {
	for curNoteID := 0; curNoteID < len(NoteList); curNoteID = curNoteID + 1 {
		curNote := NoteList[curNoteID]
		if curNote.NoteID == noteID {
			return curNote, nil
		}
	}
	return Note{}, fmt.Errorf("notes: no note exists with ID %s", noteID)
}

func getAllNotes(userID string) ([]Note, error) {
	return NoteList, nil
}
