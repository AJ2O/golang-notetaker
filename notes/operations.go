package notes

import (
	"errors"
	"net/http"

	"github.com/segmentio/ksuid"
)

// CreateNote creates a new note to be stored in the user's account.
func CreateNote(w http.ResponseWriter, r *http.Request) error {
	// parse submitted form values
	newNote := r.FormValue("note")

	// write to queue
	NoteList = append(NoteList, Note{NoteID: ksuid.New().String(), Message: newNote})

	return nil
}

// ReadNote returns the stored note with the given ID.
func ReadNote(w http.ResponseWriter, r *http.Request, noteID string) (Note, error) {
	// parse submitted form values
	userID := r.FormValue("userID")

	// get note from DB
	note, err := getNote(userID, noteID)
	if err != nil {
		return Note{}, err
	}
	return note, nil
}

// ReadAllNotes returns all the notes stored in the user's account.
func ReadAllNotes(w http.ResponseWriter, r *http.Request) ([]Note, error) {
	// parse submitted form values
	userID := r.FormValue("userID")
	return getAllNotes(userID)
}

// UpdateNote modifies the stored note with newly submitted contents.
func UpdateNote(w http.ResponseWriter, r *http.Request, noteID string) error {
	// parse submitted form values
	//userID := r.FormValue("userID")
	updateNote := r.FormValue("note")

	/*/ get note from DB
	oldNote, err := getNote(userID, noteID)
	if err != nil {
		return err
	}*/

	// update note
	for curNoteID := 0; curNoteID < len(NoteList); curNoteID = curNoteID + 1 {
		curNote := NoteList[curNoteID]
		if curNote.NoteID == noteID {
			NoteList[curNoteID].Message = updateNote
			return nil
		}
	}
	return nil
}

// DeleteNote removes the stored note from the user's account.
func DeleteNote(w http.ResponseWriter, r *http.Request, noteID string) error {
	// parse submitted form values
	userID := r.FormValue("userID")

	// get note from DB
	_, err := getNote(userID, noteID)
	if err != nil {
		return err
	}

	// delete from DB
	allNotes, _ := getAllNotes(userID)
	for curNoteID := 0; curNoteID < len(allNotes); curNoteID = curNoteID + 1 {
		curNote := allNotes[curNoteID]
		if curNote.NoteID == noteID {
			if curNoteID == len(allNotes)-1 {
				NoteList = NoteList[:curNoteID]
			} else {
				NoteList = append(NoteList[:curNoteID], NoteList[curNoteID+1:]...)
			}
			return nil
		}
	}
	return errors.New("This note doesn't exist!")
}
