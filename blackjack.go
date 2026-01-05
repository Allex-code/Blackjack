package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type Card struct {
	value int
	suit  int
}

func (card Card) createCardString() string {
	var val string
	var st string

	switch card.value {
	case 11:
		val = "J"
	case 12:
		val = "Q"
	case 13:
		val = "K"
	case 14:
		val = "A"
	case 1:
		val = "A"
	default:
		val = fmt.Sprintf("%d", card.value)
	}

	switch card.suit {
	case 0:
		st = ""
	case 1:
		st = "󰣏"
	case 2:
		st = "󰣑"
	case 3:
		st = "󰣎"
	}
	return val + st
}

func (card Card) cardValue() int {
	var cardVal int

	switch card.value {
	case 11, 12, 13:
		cardVal = 10
	case 1:
		cardVal = 11
	case 14:
		cardVal = 1
	default:
		cardVal = card.value
	}

	return cardVal
}

type Deck struct {
	Cards []Card
}

func (deck *Deck) createDeck() {
	deck.Cards = []Card{}
	for i := int(0); i <= 12; i++ {
		for j := int(0); j <= 3; j++ {
			newCard := Card{}
			newCard.value = i + 1
			newCard.suit = j
			deck.Cards = append(deck.Cards, newCard)
		}
	}
}

func (deck *Deck) printableDeck() []string {
	printedDeck := []string{}
	for _, card := range deck.Cards {
		printedDeck = append(printedDeck, card.createCardString())
	}
	return printedDeck
}

func (deck *Deck) deckLenght() int {
	return len(deck.Cards)
}

func (deck *Deck) shuffle() {
	rand.Shuffle(deck.deckLenght(), func(i, j int) {
		tmp := deck.Cards[i]
		deck.Cards[i] = deck.Cards[j]
		deck.Cards[j] = tmp
	})
}

type PlayerHands struct {
	playerHand []Card
	bet        float64
	blackjack  bool
	bust       bool
}

type Game struct {
	deck          Deck
	playerHands   []PlayerHands
	dealerCards   []Card
	dealer2ndCard Card
	balance       float64
	bet           float64
	splitedGame   bool
	blackjack     bool
	splitGames    []Game
}

// Game begins
func (game *Game) deal() {
	game.dealPlayerHand(0)
	game.dealDealer()
	game.dealPlayerHand(0)
	game.storeDealer2nd()
}

func (game *Game) dealPlayerHand(h int) {
	game.playerHands[h].playerHand = append(
		game.playerHands[h].playerHand,
		game.deck.Cards[len(game.deck.Cards)-1],
	)
	game.deck.Cards = game.deck.Cards[:len(game.deck.Cards)-1]
}

func (game *Game) dealDealer() {
	game.dealerCards = append(game.dealerCards, game.deck.Cards[len(game.deck.Cards)-1])
	game.deck.Cards = game.deck.Cards[:len(game.deck.Cards)-1]
}

func (game *Game) storeDealer2nd() {
	game.dealer2ndCard = game.deck.Cards[len(game.deck.Cards)-1]
	game.deck.Cards = game.deck.Cards[:len(game.deck.Cards)-1]
}

func (game *Game) dealerAppend() {
	game.dealerCards = append(game.dealerCards, game.dealer2ndCard)
}

type Move int

const (
	Hit Move = iota + 1
	Stand
	DoubleDown
	Split
)

func (game *Game) nextMove(h int) Move {
	doubleDown, split := false, false

	if len(game.playerHands[h].playerHand) != 2 {
		fmt.Println("HIT	STAND")
	} else if len(game.playerHands[h].playerHand) == 2 && !(game.playerHands[h].playerHand[0].cardValue() == game.playerHands[h].playerHand[1].cardValue()) {
		fmt.Println("HIT	STAND	 DOUBLE_DOWN")
		doubleDown = true
	} else if (len(game.playerHands[h].playerHand) == 2) && (game.playerHands[h].playerHand[0].cardValue() == game.playerHands[h].playerHand[1].cardValue()) {
		fmt.Println("HIT	STAND	 DOUBLE_DOWN	SPLIT")
		split = true
	}

	move := bufio.NewReader(os.Stdin)
	strMove, err := move.ReadString('\n')
	if err != nil {
		fmt.Println("Bad move, try again!", err)
	}
	trimmedMove := strings.TrimSuffix(strMove, "\n")
	choice, err := strconv.Atoi(trimmedMove)
	if err != nil {
		fmt.Println("Error while parsing input to int", err)
	}

	if !doubleDown && !split && (0 >= choice || choice >= 3) {
		fmt.Println("Not valid choice !")
		game.move(h)
		return Move(choice)
	} else if doubleDown && (0 >= choice || choice >= 4) {
		fmt.Println("Not valid choice !")
		game.move(h)
		return Move(choice)
	} else if split && (0 >= choice || choice >= 5) {
		fmt.Println("Not valid choice !")
		game.move(h)
		return Move(choice)
	}

	return Move(choice)
}

func (game *Game) dealerSoftness() {
	dealer := game.dealerScore()
	for i, card := range game.dealerCards {
		if (dealer > 21) && (card.cardValue() == 11) {
			game.dealerCards[i].value = 14
			return
		}
	}
}

func (game *Game) playerHandSoftness(h int) {
	player := game.playerHandScore(h)
	for i, card := range game.playerHands[h].playerHand {
		if (player > 21) && (card.cardValue() == 11) {
			game.playerHands[h].playerHand[i].value = 14
			return
		}
	}
}

func (game *Game) move(h int) {
	if game.deck.deckLenght() < 1 {
		game.deck.createDeck()
		game.deck.shuffle()
	}
	fmt.Printf("\nDeck: %d\n\n", game.deck.deckLenght())

	game.blackjack = false
	player := game.playerHandScore(h)

	if len(game.playerHands[0].playerHand) == 2 && !game.splitedGame {
		if player == 21 {
			game.blackjack = true
			return
		} else if player == 22 {
			game.playerHands[h].playerHand[0].value = 14
			player = game.playerHandScore(0)
		}
	}

	choice := game.nextMove(h)

	switch choice {
	case 1:
		fmt.Println("\nHit !\n")
		game.dealPlayerHand(h)
		game.playerHandSoftness(h)
		player = game.playerHandScore(h)
		if player < 21 {
			game.prettyPrintPlayerHand(h)
			game.prettyPrintDealerHand()

			game.move(h)
			return

		} else if player > 21 {
			game.prettyPrintPlayerHand(h)
			game.prettyPrintDealerHand()
			fmt.Printf("\n\nBUST !\n")
			game.balance -= game.playerHands[h].bet

			fmt.Printf("Balance: %.2f $", game.balance)
			fmt.Printf("\nDeck: %d\n\n", game.deck.deckLenght())
			game.playerHands[h].bust = true
			return

		} else if player == 21 {
			game.prettyPrintPlayerHand(h)
			game.prettyPrintDealerHand()
			fmt.Println()

			break
		}

	case 2:
		fmt.Println("\nStand !\n")
		if len(game.playerHands) > 1 {
			return
		}
	}
	if len(game.playerHands[h].playerHand) == 2 {
		switch choice {
		case 3:
			fmt.Println("\nDouble Down !\n")
			game.playerHands[h].bet = 2 * game.bet
			game.dealPlayerHand(h)
			game.playerHandSoftness(h)
			player := game.playerHandScore(h)
			if player > 21 {
				game.playerHandSoftness(h)
				player = game.playerHandScore(h)
				if player > 21 {
					player = game.playerHandScore(h)
					game.prettyPrintPlayerHand(h)
					game.prettyPrintDealerHand()
					fmt.Println("\nBUST !")
					game.playerHands[h].bust = true
					game.balance -= game.playerHands[h].bet
					return
				}
			}
			game.prettyPrintPlayerHand(h)
			game.prettyPrintDealerHand()
			fmt.Printf("\nDeck: %d\n\n", game.deck.deckLenght())

			if len(game.playerHands) > 1 {
				return
			}
		}
	}
	if true {
		switch choice {
		case 4:
			if len(game.playerHands[h].playerHand) == 2 {
				newHand := PlayerHands{
					playerHand: []Card{game.playerHands[0].playerHand[1]},
					bet:        game.bet,
					blackjack:  false,
				}
				game.playerHands = append(game.playerHands, newHand)
				game.playerHands[0].playerHand = game.playerHands[0].playerHand[:1]
			}

			for i := range game.playerHands {
				if game.deck.deckLenght() < 1 {
					game.deck.createDeck()
					game.deck.shuffle()
				}
				game.playerHands[i].bet = game.bet
				game.playerHands[i].blackjack = false

				if len(game.playerHands[i].playerHand) == 1 {
					game.dealPlayerHand(i)
				}
				fmt.Printf("\nHand: %d\n", i+1)
				fmt.Printf("Bet:\t\t %.2f $\n", game.playerHands[i].bet)
				game.prettyPrintPlayerHand(i)
				game.prettyPrintDealerHand()
				fmt.Printf("\nDealer 2nd: %s\n", game.dealer2ndCard.createCardString())
				player = game.playerHandScore(i)
				if player == 21 {
					continue
				}
				game.move(i)
			}

			return
		}
	}
	// If Not Mulitple hands jest evalute
	// otherwise return to the splited hands loop
	if len(game.playerHands) > 1 {
		return
	}
}

func (game *Game) filterBust() {
	filtered := []PlayerHands{}
	for _, g := range game.playerHands {
		if !g.bust {
			filtered = append(filtered, g)
		}
	}
	game.playerHands = filtered
}

func (game *Game) evaluate() {
	if len(game.playerHands) == 0 {
		return
	}

	if game.blackjack {
		game.dealerAppend()
		player := game.playerHandScore(0)
		dealer := game.dealerScore()
		game.prettyPrintPlayerHand(0)
		game.prettyPrintDealerHand()
		fmt.Printf("\nDeck: %d\n\n", game.deck.deckLenght())
		if player > dealer {
			game.balance += game.playerHands[0].bet * 1.5
			fmt.Printf("BLACKJACK !")
			return
		} else if player == dealer {
			fmt.Println("PUSH !")
			return
		}
	}

	game.dealerAppend()
	game.dealerSoftness()
	dealer := game.dealerScore()
	game.prettyPrintPlayerHand(0)
	game.prettyPrintDealerHand()
	fmt.Printf("\nDeck: %d\n\n", game.deck.deckLenght())
	for dealer < 17 {
		game.dealDealer()
		game.dealerSoftness()
		dealer = game.dealerScore()
		game.prettyPrintPlayerHand(0)
		game.prettyPrintDealerHand()
		fmt.Printf("\nDeck: %d\n", game.deck.deckLenght())
	}

	for h := range game.playerHands {
		fmt.Printf("\n\nHand: %d\n", h+1)
		player := game.playerHandScore(h)
		dealer = game.dealerScore()
		if player > dealer || dealer > 21 {
			game.balance += game.playerHands[h].bet
			fmt.Println("WIN !")
			return
		} else if player == dealer {
			fmt.Println("PUSH !")
			return
		} else if player < dealer {
			game.balance -= game.playerHands[h].bet
			fmt.Printf("BUST !")
		}
	}
}

func (game *Game) playerHandScore(h int) int {
	var playerSum int
	for _, card := range game.playerHands[h].playerHand {
		playerSum += card.cardValue()
	}
	return playerSum
}

func (game *Game) dealerScore() int {
	var dealerSum int
	for _, card := range game.dealerCards {
		dealerSum += card.cardValue()
	}
	return dealerSum
}

func (game *Game) printPlayerHand(h int) []string {
	printedHand := []string{}

	for _, card := range game.playerHands[h].playerHand {
		printedHand = append(printedHand, card.createCardString())
	}

	return printedHand
}

func (game *Game) printDealerHand() []string {
	printedHand := []string{}
	for _, card := range game.dealerCards {
		printedHand = append(printedHand, card.createCardString())
	}
	return printedHand
}

func (game *Game) prettyPrintPlayerHand(h int) {
	fmt.Printf("%s\n", strings.Join(game.printPlayerHand(h), " "))
}

func (game *Game) prettyPrintDealerHand() {
	fmt.Printf("%s", strings.Join(game.printDealerHand(), " "))
}

func (game *Game) balanc() float64 {
	fmt.Println("Your buy-in: ")
	readBalance := bufio.NewReader(os.Stdin)
	strBalance, err := readBalance.ReadString('\n')
	if err != nil {
		fmt.Println("Invalid Balance, either too HIGH, or invalid format !", err)
	}

	trimStrBalance := strings.TrimSuffix(strBalance, "\n")

	balance, err := strconv.ParseFloat(trimStrBalance, 64)
	if err != nil {
		fmt.Println("Error parsing string to float, check it out!")
	}
	return balance
}

func (game *Game) enterBet() float64 {
	fmt.Println("Enter your bet: ")
	reader := bufio.NewReader(os.Stdin)
	betString, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error while reading bet: ", err)
	}

	trimmedBetString := strings.TrimSuffix(betString, "\n")
	betFloat, err := strconv.ParseFloat(trimmedBetString, 64)
	if err != nil {
		fmt.Println("Error while parsing to float !", err)
	}

	return betFloat
}

func playAgain() bool {
	playAgain := false
	fmt.Printf("(1)New Game		(2)Cash Out\n")
	play := bufio.NewReader(os.Stdin)
	playStr, err := play.ReadString('\n')
	if err != nil {
		fmt.Println("Thats not what you can do!", err)
	}
	trimmPlayStr := strings.TrimSuffix(playStr, "\n")
	choice, err := strconv.Atoi(trimmPlayStr)
	if err != nil {
		fmt.Println("Not a number!", err)
	}
	if choice == 1 {
		playAgain = true
	} else if choice == 2 {
		playAgain = false
	}

	return playAgain
}

func (game *Game) newGame() {
	newGame := &Game{}
	game.splitedGame = false
	newGame.balance = game.balance
	//	newGame.bet = newGame.enterBet()
	newGame.bet = 200
	newGame.playerHands = []PlayerHands{{bet: newGame.bet}}
	fmt.Printf("Bet:\t\t %.2f $\n\n", newGame.bet)
	// TODO: len(newGame.dec)
	if game.deck.deckLenght() < 5 {
		game.deck.createDeck()
		game.deck.shuffle()
		newGame.deck = game.deck
	}
	newGame.deck = game.deck
	newGame.deal()

	newGame.prettyPrintPlayerHand(0)
	newGame.prettyPrintDealerHand()
	newGame.move(0)
	newGame.filterBust()
	newGame.evaluate()
	game.balance = newGame.balance
	game.deck = newGame.deck
	fmt.Printf("\nTotal:\t %.2f $\n", game.balance)
	play := playAgain()
	if play == true {
		game.newGame()
	} else {
		return
	}
}

func main() {
	game := &Game{}
	//	game.balance = game.balanc()
	game.balance = 10000
	fmt.Printf("Balance:\t %.2f $\n", game.balance)
	game.newGame()
}
