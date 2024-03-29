package go_wordle_solver

import (
	"fmt"

	"golang.org/x/exp/slices"
)

// Indicates if a letter is known to be in a given location or not.
type locatedLetterState uint8

const (
	llsUnknown locatedLetterState = iota
	llsHere
	llsNotHere
)

func (lls locatedLetterState) String() string {
	switch lls {
	case llsUnknown:
		return "unknown"
	case llsHere:
		return "here"
	case llsNotHere:
		return "not here"
	}
	return "invalid value"
}

// Indicates information about a letter that is in the word.
type presentLetter struct {
	// If known, the letter must appear exactly this many times in the word.
	maybeRequiredCount Optional[uint8]
	// The minimum number of times this letter must appear in the word.
	minCount uint8
	// The number of locations we know the letter must appear.
	numHere uint8
	// The number of locations we know the letter must not appear.
	numNotHere uint8
	// The status of the letter at each location in the word.
	locatedState []locatedLetterState
}

// Constructs a [presentLetter] for use with words of the given length.
func newPresentLetter(wordLength uint8) *presentLetter {
	states := make([]locatedLetterState, wordLength)
	for i := range states {
		states[i] = llsUnknown
	}
	pl := new(presentLetter)
	pl.minCount = 1
	pl.locatedState = states
	return pl
}

// Returns whether the letter must be in, or not in, the given location, or if that is not yet
// known.
func (self *presentLetter) state(index uint8) locatedLetterState {
	return self.locatedState[index]
}

// Sets that this letter must be at the given index.
//
// If the required count for this letter is known, then this may fill any remaining [llsUnknown]
// locations with either [llsHere] or [llsNotHere] accordingly.
//
// This returns an error if this letter is already known not to be at the given index.
func (self *presentLetter) setMustBeAt(index uint8) error {
	previous := self.locatedState[index]
	switch previous {
	case llsHere:
		return nil
	case llsNotHere:
		return fmt.Errorf("Can't set letter to %s at index %v since it's already marked as %s.", llsHere, index, previous)
	}
	self.locatedState[index] = llsHere
	self.numHere += 1
	if self.numHere > self.minCount {
		self.minCount = self.numHere
	}
	if self.maybeRequiredCount.HasValue() {
		count := self.maybeRequiredCount.Value()
		if self.numHere == count {
			// If the count has been met, then this letter doesn't appear anywhere else.
			self.setUnknownsTo(llsNotHere)
		} else if (uint8(len(self.locatedState)) - self.numNotHere) == count {
			// If the letter must be in all possible remaining spaces, set them to here.
			self.setUnknownsTo(llsHere)
		}
	} else {
		// Set the max count if all states are known to prevent errors.
		// Note that there is no need to update any unknowns in this case, as there are no
		// unknowns left.
		self.setRequiredCountIfFull()
	}
	return nil
}

// Sets that this letter must not be at the given index.
//
// If setting this leaves only as many [llsHere] and [llsUnknown] locations as the value of
// [presentLetter.minCount], then this sets the [llsUnknown] locations to [llsHere].
//
// This returns an error if this letter is already known to be at the given index.
func (self *presentLetter) setMustNotBeAt(index uint8) error {
	previous := self.locatedState[index]
	switch previous {
	case llsNotHere:
		return nil
	case llsHere:
		return fmt.Errorf("Can't set letter to %s at index %v since it's already marked as %s.", llsNotHere, index, previous)
	}
	self.locatedState[index] = llsNotHere
	self.numNotHere += 1
	maxPossibleHere := uint8(len(self.locatedState)) - self.numNotHere
	if maxPossibleHere == self.minCount {
		// If the letter must be in all possible remaining spaces, set them to `Here`.
		self.maybeRequiredCount = OptionalOf(self.minCount)
		if self.numHere < self.minCount {
			self.setUnknownsTo(llsHere)
		}
	}
	return nil
}

// Sets the maximum number of times this letter can appear in the word.
//
// Returns an error if the required count is already set to a different value, or if the
// [presentLetter.minCount] is known to be higher than the provided value.
func (self *presentLetter) setRequiredCount(count uint8) error {
	if self.maybeRequiredCount.HasValue() {
		if self.maybeRequiredCount.Value() != count {
			return fmt.Errorf("Can't set required count to %v since it's already %v.", count, self.maybeRequiredCount.Value())
		} else {
			return nil
		}
	}
	if self.minCount > count {
		return fmt.Errorf("Can't set required count to %v since that would be less than the minimum count (%v).", count, self.minCount)
	}
	self.minCount = count
	maxPossibleNumHere := uint8(len(self.locatedState)) - self.numNotHere
	if maxPossibleNumHere < count {
		return fmt.Errorf("Can't set required count to %v since there aren't enough possible spaces (only %v).", count, maxPossibleNumHere)
	}
	self.maybeRequiredCount = OptionalOf(count)
	if self.numHere == count {
		self.setUnknownsTo(llsNotHere)
	} else if maxPossibleNumHere == count {
		self.setUnknownsTo(llsHere)
	}
	return nil
}

// If count is higher than the current min-count, this bumps it up to the provided value and
// modifies the known data as needed.
//
// Returns an error if it would be impossible for count locations to be marked [llsHere] given what
// is already known about the word.
func (self *presentLetter) possiblyBumpMinCount(count uint8) error {
	if self.minCount >= count {
		return nil
	}

	self.minCount = count
	maxPossibleNumHere := uint8(len(self.locatedState)) - self.numNotHere
	if maxPossibleNumHere < count {
		return fmt.Errorf("Can't set min count to %v when there are only %v possible locations.", count, maxPossibleNumHere)
	} else if maxPossibleNumHere == count && self.numHere < count {
		// If all possible unknowns must be here, set them.
		self.setUnknownsTo(llsHere)
		self.maybeRequiredCount = OptionalOf(count)
	}
	return nil
}

func (self *presentLetter) setUnknownsTo(newState locatedLetterState) {
	if newState == llsUnknown {
		// Nothing to do.
		return
	}

	countToUpdate := &self.numHere
	if newState == llsNotHere {
		countToUpdate = &self.numNotHere
	}
	for i, state := range self.locatedState {
		if state == llsUnknown {
			self.locatedState[i] = newState
			*countToUpdate += 1
		}
	}
}

func (self *presentLetter) setRequiredCountIfFull() {
	if self.numHere+self.numNotHere == uint8(len(self.locatedState)) {
		self.maybeRequiredCount = OptionalOf(self.numHere)
	}
}

// LetterRestriction indicates the known restrictions that apply to a letter at a given location.
//
// See [WordRestrictions].
type LetterRestriction uint8

const (
	// The letter restriction is unknown.
	LetterRestrictionUnknown LetterRestriction = iota
	// The letter goes here.
	LetterRestrictionHere
	// The letter is in the word and might be here.
	LetterRestrictionPresentMaybeHere
	// The letter is in the word but not here.
	LetterRestrictionPresentNotHere
	// The letter is not in the word.
	LetterRestrictionNotPresent
)

func (lr LetterRestriction) String() string {
	switch lr {
	case LetterRestrictionHere:
		return "here"
	case LetterRestrictionPresentMaybeHere:
		return "present maybe here"
	case LetterRestrictionPresentNotHere:
		return "present not here"
	case LetterRestrictionNotPresent:
		return "not present"
	case LetterRestrictionUnknown:
		return "unknown"
	}
	return "invalid value"
}

// WordRestrictions defines letter restrictions that a word must adhere to, such as "the first
// letter of the word must be 'a'".
//
// Restrictions are derived from [GuessResult]s.
type WordRestrictions struct {
	wordLength        uint8
	presentLetters    map[rune]*presentLetter
	notPresentLetters []rune
}

// Creates a [WordRestrictions] object for the given word length with all letters unknown.
func InitWordRestrictions(wordLength uint8) WordRestrictions {
	return WordRestrictions{
		wordLength,
		make(map[rune]*presentLetter, wordLength),
		make([]rune, 0, 13),
	}
}

// WordRestrictionsFromResult returns the restrictions imposed by the given result.
func WordRestrictionsFromResult(result *GuessResult) (WordRestrictions, error) {
	restrictions := InitWordRestrictions(uint8(result.Guess.Len()))
	err := restrictions.Update(result)
	return restrictions, err
}

// Update adds restrictions arising from the given result.
//
// Returns an error if the result is incompatible with the existing restrictions.
func (self *WordRestrictions) Update(guessResult *GuessResult) error {
	var err error
	guessLen := guessResult.Guess.Len()
	for i := 0; i < guessLen; i++ {
		letter := guessResult.Guess.At(i)
		switch guessResult.Results[i] {
		case LetterResultCorrect:
			err = self.setLetterHere(letter, uint8(i), guessResult)
		case LetterResultPresentNotHere:
			err = self.setLetterPresentNotHere(letter, uint8(i), guessResult)
		case LetterResultNotPresent:
			err = self.setLetterNotPresent(letter, uint8(i), guessResult)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// IsSatisfiedBy returns true iff the given word satisfies these restrictions.
func (self *WordRestrictions) IsSatisfiedBy(word Word) bool {
	wordLen := word.Len()
	return wordLen == int(self.wordLength) &&
		allPairs(self.presentLetters, func(letter rune, presence *presentLetter) bool {
			countFound := uint8(0)
			for i := 0; i < wordLen; i++ {
				wordLetter := word.At(i)
				if wordLetter == letter {
					countFound += 1
					if presence.state(uint8(i)) == llsNotHere {
						return false
					}
				} else if presence.state(uint8(i)) == llsHere {
					return false
				}
			}
			if presence.maybeRequiredCount.HasValue() {
				return countFound == presence.maybeRequiredCount.Value()
			}
			return countFound >= presence.minCount
		}) &&
		word.AllLetters(func(letter rune) bool {
			return !slices.Contains(self.notPresentLetters, letter)
		})
}

// IsStateKnown returns true iff the exact state of the given letter at the given location is
// already known.
func (self *WordRestrictions) IsStateKnown(letter rune, location uint8) bool {
	if presence, isPresent := self.presentLetters[letter]; isPresent {
		return presence.state(location) != llsUnknown
	}
	return slices.Contains(self.notPresentLetters, letter)
}

// State returns the current known state of this letter at the given location, either:
//
//  * [LetterRestrictionUnknown] -> Nothing is known about the letter.
//  * [LetterRestrictionNotPresent] -> The letter is not in the word.
//  * [LetterRestrictionPresentNotHere] -> The letter is present but not here.
//  * [LetterRestrictionPresentMaybeHere] -> The letter is present, but we don't know if it's here or not.
//  * [LetterRestrictionHere] -> The letter goes here.
func (self *WordRestrictions) State(letter rune, location uint8) LetterRestriction {
	if presence, isPresent := self.presentLetters[letter]; isPresent {
		switch presence.state(location) {
		case llsHere:
			return LetterRestrictionHere
		case llsNotHere:
			return LetterRestrictionPresentNotHere
		default:
			return LetterRestrictionPresentMaybeHere
		}
	}
	if slices.Contains(self.notPresentLetters, letter) {
		return LetterRestrictionNotPresent
	}
	return LetterRestrictionUnknown
}

func (self *WordRestrictions) setLetterHere(
	letter rune,
	location uint8,
	result *GuessResult,
) error {
	presence, isPresent := self.presentLetters[letter]
	if !isPresent {
		presence = newPresentLetter(self.wordLength)
		self.presentLetters[letter] = presence
	}
	err := presence.setMustBeAt(location)
	if err != nil {
		return err
	}
	numTimesPresent := countNumTimesInGuess(letter, result)
	// Remove from the not present letters if it was present. This could happen if the guess
	// included the letter in two places, but the correct word only included it in the latter
	// place.
	if letterIndex := slices.Index(self.notPresentLetters, letter); letterIndex >= 0 {
		self.notPresentLetters = slices.Delete(self.notPresentLetters, letterIndex, letterIndex+1)
		// If the letter is present, but another location was marked `NotPresent`, then it means it's
		// only in the word as many times as it was given a `Correct` or  `PresentNotHere`
		// hint.
		err = presence.setRequiredCount(numTimesPresent)
	} else {
		err = presence.possiblyBumpMinCount(numTimesPresent)
	}
	if err != nil {
		return err
	}

	// Remove this location possibility from other present letters.
	for otherLetter, otherPresence := range self.presentLetters {
		if otherLetter == letter {
			continue
		}

		err = otherPresence.setMustNotBeAt(location)
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *WordRestrictions) setLetterPresentNotHere(
	letter rune,
	location uint8,
	result *GuessResult,
) error {
	presence, isPresent := self.presentLetters[letter]
	if !isPresent {
		presence = newPresentLetter(self.wordLength)
		self.presentLetters[letter] = presence
	}
	err := presence.setMustNotBeAt(location)
	if err != nil {
		return err
	}
	numTimesPresent := countNumTimesInGuess(letter, result)
	return presence.possiblyBumpMinCount(numTimesPresent)
}

func (self *WordRestrictions) setLetterNotPresent(
	letter rune,
	location uint8,
	result *GuessResult,
) error {
	numTimesPresent := countNumTimesInGuess(letter, result)
	if presence, isPresent := self.presentLetters[letter]; isPresent {
		if presence.state(location) == llsHere {
			return fmt.Errorf("Can't mark the letter %c as not present at %v since it's already marked as present here.", letter, location)
		}
		return presence.setRequiredCount(numTimesPresent)
	} else if numTimesPresent > 0 {
		pl := newPresentLetter(uint8(result.Guess.Len()))
		self.presentLetters[letter] = pl
		return pl.setRequiredCount(numTimesPresent)
	}
	if !slices.Contains(self.notPresentLetters, letter) {
		self.notPresentLetters = append(self.notPresentLetters, letter)
	}
	return nil
}

func countNumTimesInGuess(letter rune, guessResult *GuessResult) uint8 {
	var sum uint8 = 0
	guessLen := guessResult.Guess.Len()
	for i := 0; i < guessLen; i++ {
		guessLetter := guessResult.Guess.At(i)
		if guessLetter == letter && guessResult.Results[i] != LetterResultNotPresent {
			sum += 1
		}
	}
	return sum
}
