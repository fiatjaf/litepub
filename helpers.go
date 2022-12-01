package litepub

func WrapCreate(note Note, createId string) (create Create[Note]) {
	return Create[Note]{
		Base: Base{
			Type: "Create",
			Id:   createId,
		},
		Actor:  note.AttributedTo,
		Object: note,
	}
}

func unwrapCreate(create Create[Note]) Note {
	return create.Object
}
