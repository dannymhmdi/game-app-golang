package entity

type Question struct {
	ID              uint
	Text            string
	PossibleAnswers []PossibleAnswer
	CorrectAnswerID uint
	Difficulty      QuestionDifficulty
	CategoryID      uint
}

// question option
type PossibleAnswer struct {
	ID     uint
	Text   string
	Choice PossibleAnswerChoice
}

type PossibleAnswerChoice uint8

func (p PossibleAnswerChoice) IsValid() bool {
	if p >= PossibleAnswerA && p <= PossibleAnswerD {
		return true
	}
	return false
}

const (
	PossibleAnswerA PossibleAnswerChoice = iota + 1
	PossibleAnswerB
	PossibleAnswerC
	PossibleAnswerD
)

type QuestionDifficulty uint8

func (d QuestionDifficulty) IsValid() bool {
	if d >= 1 && d <= 3 {
		return true
	}
	return false
}

const (
	questionDifficultyEasy QuestionDifficulty = iota + 1
	questionDifficultyMedium
	questionDifficultyHard
)

type Validator interface {
	IsValid() bool
}
