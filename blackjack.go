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

type Game struct {
	deck          Deck
	playerCards   []Card
	dealerCards   []Card
	dealer2ndCard Card
	balance       float64
	bet           float64
}

// Game begins
func (game *Game) deal() {
	game.dealPlayer()
	game.dealDealer()
	game.dealPlayer()
	game.storeDealer2nd()
}

func (game *Game) dealPlayer() {
	game.playerCards = append(game.playerCards, game.deck.Cards[len(game.deck.Cards)-1])
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

func (game *Game) nextMove() Move {
	doubleDown, split := false, false

	game.gameState()
	if len(game.playerCards) != 2 {
		fmt.Println("HIT	STAND")
	} else if (len(game.playerCards) == 2) && (game.playerCards[0].value == game.playerCards[1].value) {
		fmt.Println("HIT	STAND	 DOUBLE_DOWN	SPLIT")
		split = true
	} else if len(game.playerCards) == 2 {
		fmt.Println("HIT	STAND	 DOUBLE_DOWN")
		doubleDown = true

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
		game.move()
		return Move(choice)
	} else if doubleDown && (0 >= choice || choice >= 4) {
		fmt.Println("Not valid choice !")
		game.move()
		return Move(choice)
	} else if split && (0 >= choice || choice >= 5) {
		fmt.Println("Not valid choice !")
		game.move()
		return Move(choice)
	}

	return Move(choice)
}

func (game *Game) dealerSoftness() int {
	dealer := game.dealerScore()
	for i, card := range game.dealerCards {
		if (dealer > 21) && (card.cardValue() == 11) {
			game.dealerCards[i].value = 14
			return dealer
		}
	}
	dealer = game.dealerScore()

	return dealer
}


func (game *Game) playerSoftness() int {
	dealer := game.playerScore()
	for i, card := range game.playerCards{
		if (dealer > 21) && (card.cardValue() == 11) {
			game.playerCards[i].value = 14
			return dealer
		}
	}
	dealer = game.playerScore()

	return dealer
}


func (game *Game) move() {
	if game.deck.deckLenght() < 1 {
		game.deck.createDeck()
		game.deck.shuffle()
	}

	player := game.playerScore()
	dealer := game.dealerScore()
	bet := game.bet

	if len(game.playerCards) == 2 {
		if player == 21 {
			game.dealerAppend()
			dealer = game.dealerScore()
			if dealer != 21 {
				bet = bet * 1.5
				game.balance += bet
				game.gameState()
				fmt.Println("BlackJack !!!!\t", bet, "$")
			} else if dealer == 21 {
				game.gameState()
				fmt.Println("Push !!!!")
			}
			return
		} else if player == 22 {
			game.playerCards[0].value = 14
			game.gameState()
			player = game.playerScore()
		}
	}

	choice := game.nextMove()

	switch choice {
	case 1:
		fmt.Println("\nHit !")
		game.dealPlayer()
		player = game.playerScore()
		if player == 21 {
			game.dealerAppend()
			dealer = game.dealerScore()
			game.gameState()
			for dealer <= 16 {
				game.dealDealer()
				game.dealerSoftness()
				dealer = game.dealerScore()
			}
			if (dealer < 21) || (dealer > 21) {
				game.gameState()
				game.balance += bet
				fmt.Println("\nWon !")
			} else if dealer == 21 {
				game.gameState()
				fmt.Println("\nPush !")
			}
			return
		} else if player < 21 {
			game.gameState()
			game.move()
			return
		} else if player > 21 {
			for i, card := range game.playerCards {
				if card.cardValue() == 11 {
					game.playerCards[i].value = 14
				}
				if game.playerScore() < 21 {
					game.gameState()
					game.move()
					return
				}
			}
			if len(game.dealerCards) == 1 {
				// game.dealerAppend()
				game.gameState()
				game.balance -= bet
				fmt.Println("\nBust !")
				return
			}
		}
	case 2:
		fmt.Println("\nStand !")
		game.dealerAppend()
		game.dealerSoftness()
		player = game.playerScore()
		dealer = game.dealerScore()
		game.gameState()

		for dealer <= 16 {
			game.dealDealer()
			game.dealerSoftness()
			dealer = game.dealerScore()
		}
		if (player > dealer) || (dealer > 21) {
			game.gameState()
			game.balance += bet
			fmt.Println("Won !")
			return
		} else if dealer > player {
			game.gameState()
			game.balance -= bet
			fmt.Println("Bust !")
			return
		} else {
			game.gameState()
			fmt.Println("Push")
			return
		}
	}
	if len(game.playerCards) == 2 {
		switch choice {
		case 3:
			fmt.Println("\nDouble Down !")
			bet = 2 * bet
			game.dealPlayer()
			player = game.playerScore()
			if player > 21 {
				for i, card := range game.playerCards {
					if card.cardValue() == 11 {
						game.playerCards[i].value = 14
						player = game.playerScore()
						break
					}
				}
			}

			game.dealerAppend()
			game.dealerSoftness()
			game.gameState()
			dealer := game.dealerScore()

			if player <= 21 {
				for dealer <= 16 {
					game.dealDealer()
					game.dealerSoftness()
					game.gameState()
					dealer = game.dealerScore()
				}
				if (dealer < player) || (dealer > 21) {
					game.balance += bet
					fmt.Println("\nWon !!!a")
				} else if dealer == player {
					fmt.Println("\nPush !!!a")
				} else if dealer > player  {
					game.balance -= bet
					fmt.Println("\nBust !!!a")
				}
				return
			} else if (player > 21) {
				game.balance -= bet
				fmt.Println("Bust !")
			}
		}
	}
	if true {
		switch choice {
		case 4:
			fmt.Println("Split !")	
			
		}
	}
}

func (game *Game) playerScore() int {
	var playerSum int
	for _, card := range game.playerCards {
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

func (game *Game) printPlayerHand(b bool) []string {
	printedHand := []string{}

	switch b {
	case true:
		for _, card := range game.playerCards {
			printedHand = append(printedHand, card.createCardString())
		}
	case false:
		for _, card := range game.dealerCards {
			printedHand = append(printedHand, card.createCardString())
		}
	}
	return printedHand
}

func (game *Game) prettyPrintHands() {
	fmt.Printf("\nPlayer:\t\t %s\n", strings.Join(game.printPlayerHand(true), " "))
	fmt.Printf("Dealer:\t\t %s\n", strings.Join(game.printPlayerHand(false), " "))
}

func (game *Game) gameState() {
	playerSum := game.playerScore()
	dealerSum := game.dealerScore()
	game.prettyPrintHands()
	fmt.Printf("Player score:\t %d\n", playerSum)
	fmt.Printf("Dealer score:\t %d\n", dealerSum)
	fmt.Println("Deck: ", game.deck.deckLenght())
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
	newGame.balance = game.balance
	//	newGame.bet = newGame.enterBet()
	newGame.bet = 200
	fmt.Printf("Bet:\t\t %.2f $\n", newGame.bet)
	// TODO: len(newGame.dec)
	if game.deck.deckLenght() < 5 {
		game.deck.createDeck()
		game.deck.shuffle()
		newGame.deck = game.deck
	}
	newGame.deck = game.deck
	newGame.deal()

	newGame.move()
	game.balance = newGame.balance
	game.deck = newGame.deck
	fmt.Printf("\nTotal :\t %.2f $\n", game.balance)
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
