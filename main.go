package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/nexidian/gocliselect"
	"strconv"

	"pokemon-fight-cli/pokemons"
	"pokemon-fight-cli/potions"
)

func pokemonTeamAlive(team map[string]map[string]interface{}) bool {
	for _, stats := range team {
		if stats["HP"].(int)>0{
			return true
		}
	}
	return false
}

func printMessage(str string) {
	fmt.Println(str)
	time.Sleep(1*time.Second)
}

func isPokemonAlive(pokemon string, player *Player) bool {
	return player.team[pokemon]["HP"].(int)>0
}

func applyDmg(attacker Player, victim Player) {
	attackerPokemon:=attacker.mainPokemon
	attackerPokemonStats:=attacker.team[attackerPokemon]

	victim_HP:=victim.team[victim.mainPokemon]["HP"].(int)

	if victim_HP-attackerPokemonStats["Dmg"].(int)<0{
		victim.team[victim.mainPokemon]["HP"]=0
	} else {
		victim.team[victim.mainPokemon]["HP"]=victim_HP-attackerPokemonStats["Dmg"].(int)
	}
}

func applyPotion(player *Player){
	pokemonHP:=player.team[player.mainPokemon]["HP"].(int)
	
	if pokemonHP+2>player.team[player.mainPokemon]["MaxHP"].(int){
		player.team[player.mainPokemon]["HP"]=player.team[player.mainPokemon]["MaxHP"].(int)
	} else {
		player.team[player.mainPokemon]["HP"]=pokemonHP+2
	}
}


//===== EN CONSTRUCTION =====
func giveEffects(attacker *Player, victim *Player) [][]string {
	//Caractéristiques de chaques effets
	victimEffects:=victim.team[victim.mainPokemon]["Effects"].([][]string)

	for _, typAttacker := range attacker.team[attacker.mainPokemon]["Type"].([]string){
		//Vérification si les types sont immunisés
		sameType:=false
		for _, typVictim := range victim.team[victim.mainPokemon]["Type"].([]string){
			if typVictim==typAttacker && typAttacker!="Psychic"{
				sameType=true
			}
		}
		if sameType {
			continue
		}

		isBurnt:=false
		isFroze:=false
		isPoison:=false
		isParalyse:=false
		isPsychic:=false
		for _, e := range victimEffects {
			if e[0]=="brûlure"{
				isBurnt=true
			} else if e[0]=="gel"{
				isFroze=true
			} else if e[0]=="poison"{
				isPoison=true
			} else if e[0]=="paralysie"{
				isParalyse=true
			} else if e[0]=="sommeil"{
				isPsychic=true
			}
		}

		if typAttacker=="Fire" && !isBurnt{
			random:=rand.Float64()
			if random<=1{
				victimEffects=append(victimEffects, []string{"brûlure", "3"})
			}
		} else if typAttacker=="Ice" && !isFroze{
			random:=rand.Float64()
			if random<=0.3{
				victimEffects=append(victimEffects, []string{"gel", "1"})
			}		
		} else if typAttacker=="Poison" && !isPoison{
			random:=rand.Float64()
			if random<=0.15{
				victimEffects=append(victimEffects, []string{"poison", "2"})
			}		
		} else if typAttacker=="Electric" && !isParalyse{
			random:=rand.Float64()
			if random<=0.13{
				victimEffects=append(victimEffects, []string{"paralysie", "4"})
			}		
		} else if typAttacker=="Psychic" && !isPsychic{
			random:=rand.Float64()
			if random<=0.08{
				victimEffects=append(victimEffects, []string{"sommeil", "2"})
			}		
		}
	}

	return victimEffects
}

func getHPBar(HP int, MaxHP int) string {
	HPBar := ""
	if HP<=0 {
		HPBar=strings.Repeat("□", 15)
	} else if MaxHP==HP{
		HPBar=strings.Repeat("■", 15)
	} else {
		HPBar=strings.Repeat("■", int(HP*15/MaxHP))+strings.Repeat("□", 15-int(HP*15/MaxHP))
	}
	return HPBar
}

func choosePokemonTeam(pokemons map[string]map[string]interface{}) map[string]map[string]interface{} {
    pokemonsCopy := make(map[string]map[string]interface{})
    for key, value := range pokemons {
        pokemonsCopy[key] = value
    }

	//Permet de générer aléatoirement la team pokémon
	team := make(map[string]map[string]interface{})
	 
	for i:=0; i<5; i++{
		choice := rand.Intn(len(pokemonsCopy))
		count := 0
		for pokemon, stats := range pokemonsCopy {
			if count==choice{
				team[pokemon]=stats
				team[pokemon]["Effects"]=[][]string{}
				delete(pokemonsCopy, pokemon)
				break
			}
			count+=1
		}
	}
	return team
}

func choosePotions(potions map[string]map[string]interface{}) [7]string {
	//Permet de générer aléatoirement l'inventaire d'objets
	items := [7]string {}

	for i:=0; i<len(items); i++{
		choice := rand.Float64()
		cumulatif := 0.0
		for p, _ := range potions{

			cumulatif += 1-float64(potions[p]["DropRate"].(int))/100
			if choice <= cumulatif{
				items[i]=p
				break
			}  
		}
	}

	return items
}

func getFirstPokemonTeam(team map[string]map[string]interface{}) string {
	for pokemon, _ := range team{
		return pokemon
	}
	return ""
}

func displayPokemons(player1 Player, player2 Player){
	//=== Player 1 ===
	p1_HP:=player1.team[player1.mainPokemon]["HP"].(int)
	p1_MaxHP:=player1.team[player1.mainPokemon]["MaxHP"].(int)
	p1_HPBar := getHPBar(p1_HP, p1_MaxHP)


	p1_Effects := player1.team[player1.mainPokemon]["Effects"].([][]string)
	p1_EffectsFormat:=""
	if len(p1_Effects)!=0{
		p1_EffectsFormat:="["
		for i := 0; i < len(p1_Effects); i++ {
			value := p1_Effects[i][0]
			if i < len(p1_Effects)-1 {
				p1_EffectsFormat+=string(value)+", "
			} else {
				p1_EffectsFormat+=string(value)
			}
			
		}
		p1_EffectsFormat+="]"
	}

	//=== Player 2 ===
	p2_HP:=player2.team[player2.mainPokemon]["HP"].(int)
	p2_MaxHP:=player2.team[player2.mainPokemon]["MaxHP"].(int)
	p2_HPBar := getHPBar(p2_HP, p2_MaxHP)

	p2_Effects := player2.team[player2.mainPokemon]["Effects"].([][]string)
	p2_EffectsFormat:=""
	if len(p2_Effects)!=0{
		p2_EffectsFormat:="["
		for i := 0; i < len(p2_Effects); i++ {
			value := p2_Effects[i][0]
			if i < len(p2_Effects)-1 {
				p2_EffectsFormat+=string(value)+", "
			} else {
				p2_EffectsFormat+=string(value)
			}
			
		}
		p2_EffectsFormat+="]"
	}

	fmt.Println("------------------------------------------------------------")
	fmt.Print("[JOUEUR 1]")
	fmt.Print("			")
	fmt.Print("[JOUEUR 2]")
	fmt.Print("\n")
	if len(player1.mainPokemon)<=4{
		fmt.Print("👾 "+player1.mainPokemon+"	")
	} else {
		fmt.Print("👾 "+player1.mainPokemon)
	}
	fmt.Print("			")
	fmt.Print("👾 "+player2.mainPokemon)
	fmt.Print("\n")
	fmt.Print(p1_HPBar+"["+strconv.Itoa(p1_HP)+"/"+strconv.Itoa(p1_MaxHP)+"] HP")
	fmt.Print("	")
	fmt.Print(p2_HPBar+"["+strconv.Itoa(p2_HP)+"/"+strconv.Itoa(p2_MaxHP)+"] HP")
	fmt.Print("\n")
	fmt.Print(p1_EffectsFormat+" "+p2_EffectsFormat)
	fmt.Print("\n\n")
}

func selectPokemon(player *Player) string {
	teamMenu := gocliselect.NewMenu("Joueur "+string(player.playerID)+" Team")

	for pokemon, stats := range player.team{
		//Types du pokémon
		types := stats["Type"].([]string)
		var typesFormat string
   		for i:=0; i<len(types)-1; i++ {
			typesFormat+=types[i]+"/"
		}
		typesFormat+=types[len(types)-1]

		//Vie
		HP:=stats["HP"].(int)
		MaxHP:=stats["MaxHP"].(int)
		HPBar := getHPBar(HP, MaxHP)
		
		//Attaque
		Dmg:=stats["Dmg"].(int)
		DmgFormat:=strconv.Itoa(Dmg)+"⚡"

		//Effets
		Effects := stats["Effects"].([][]string)
		EffectsFormat:=""
		for i := 0; i < len(Effects); i++ {
			value := Effects[i][0]
			if i < len(Effects)-1 {
				EffectsFormat+=string(value)+", "
			} else {
				EffectsFormat+=string(value)
			}
			
		}
		if len(Effects)>0 {
			EffectsFormat="["+EffectsFormat+"]"
		}
	

		if pokemon==player.mainPokemon{
			if !isPokemonAlive(pokemon, player){
				teamMenu.AddItem("​☠ ★ "+pokemon+" ["+typesFormat+"] "+HPBar+"["+strconv.Itoa(HP)+"/"+strconv.Itoa(MaxHP)+"] "+DmgFormat+" "+EffectsFormat, pokemon)	
			} else {
				teamMenu.AddItem("★ "+pokemon+" ["+typesFormat+"] "+HPBar+"["+strconv.Itoa(HP)+"/"+strconv.Itoa(MaxHP)+"] "+DmgFormat+" "+EffectsFormat, pokemon)
			}
		} else if !isPokemonAlive(pokemon, player){
			teamMenu.AddItem("​☠ "+pokemon+" ["+typesFormat+"] "+HPBar+"["+strconv.Itoa(HP)+"/"+strconv.Itoa(MaxHP)+"] "+DmgFormat+" "+EffectsFormat, pokemon)	
		} else {
			teamMenu.AddItem(pokemon+" ["+typesFormat+"] "+HPBar+"["+strconv.Itoa(HP)+"/"+strconv.Itoa(MaxHP)+"] "+DmgFormat+" "+EffectsFormat, pokemon)
		}
	}

	teamMenu.AddItem("◀ Retour", "back")
	choice := teamMenu.Display()

	return choice
}



//Struct Player
type Player struct {
	team map[string]map[string]interface{}
	items [7]string
	mainPokemon string
	playerID rune
}

func main() {
	//Création de l'inventaire et de la team pour chaque joueur
	player1:=Player{choosePokemonTeam(pokemons.Pokemons), choosePotions(potions.Potions), "", '1'}
	player2:=Player{choosePokemonTeam(pokemons.Pokemons), choosePotions(potions.Potions), "", '2'}

	player1.mainPokemon=getFirstPokemonTeam(player1.team)
	player2.mainPokemon=getFirstPokemonTeam(player2.team)

	currentPlayer:=&player1

	roundFinish:=true

	fmt.Println("\n██████╗░░█████╗░██╗░░██╗███████╗███╗░░░███╗░█████╗░███╗░░██╗\n██╔══██╗██╔══██╗██║░██╔╝██╔════╝████╗░████║██╔══██╗████╗░██║\n██████╔╝██║░░██║█████═╝░█████╗░░██╔████╔██║██║░░██║██╔██╗██║\n██╔═══╝░██║░░██║██╔═██╗░██╔══╝░░██║╚██╔╝██║██║░░██║██║╚████║\n██║░░░░░╚█████╔╝██║░╚██╗███████╗██║░╚═╝░██║╚█████╔╝██║░╚███║\n╚═╝░░░░░░╚════╝░╚═╝░░╚═╝╚══════╝╚═╝░░░░░╚═╝░╚════╝░╚═╝░░╚══╝\n")


	//Jeu
	for pokemonTeamAlive(player1.team) && pokemonTeamAlive(player2.team){
		roundFinish=false

		//Affichage de jeu en CLI
		displayPokemons(player1, player2)

		if !isPokemonAlive(currentPlayer.mainPokemon, currentPlayer) {
			pokemonSelected:=selectPokemon(currentPlayer)
			if pokemonSelected!="back"{
				if !isPokemonAlive(pokemonSelected, currentPlayer){
					printMessage("[Joueur "+string(currentPlayer.playerID)+"] "+currentPlayer.mainPokemon+" n'est plus en état de se battre!")
					continue
				}
				currentPlayer.mainPokemon=pokemonSelected
				printMessage("[Joueur "+string(currentPlayer.playerID)+"] "+currentPlayer.mainPokemon+" en avant!")
				displayPokemons(player1, player2)
			} else if pokemonSelected=="back" {
				printMessage("[Joueur "+string(currentPlayer.playerID)+"] Vous devez choisir un pokémon à envoyer!")
				continue
			}
		}

		selectMenu := gocliselect.NewMenu("Joueur "+string(currentPlayer.playerID)+" Actions")

		selectMenu.AddItem("Attaquer", "attack")
		selectMenu.AddItem("Utiliser une potion", "items")
		selectMenu.AddItem("Changer de pokémon", "change")
		selectMenu.AddItem("Quitter", "quit")
		choice := selectMenu.Display()

		//Attaque des pokémons
		if choice=="attack"{
			if currentPlayer.playerID=='1'{
				applyDmg(player1, player2)
				player2.team[player2.mainPokemon]["Effects"]=giveEffects(&player1, &player2)
			} else {
				applyDmg(player2, player1)
				player1.team[player1.mainPokemon]["Effects"]=giveEffects(&player2, &player1)
			}
			roundFinish=true
			printMessage("[Joueur "+string(currentPlayer.playerID)+"] "+currentPlayer.mainPokemon+" attaque! [-"+strconv.Itoa(currentPlayer.team[currentPlayer.mainPokemon]["Dmg"].(int))+"HP]")
		}

		//Attaque des pokémons
		if choice=="items"{
			if currentPlayer.team[currentPlayer.mainPokemon]["HP"].(int)==currentPlayer.team[currentPlayer.mainPokemon]["MaxHP"].(int){
				printMessage("[Joueur "+string(currentPlayer.playerID)+"] Impossible d'utiliser une potion !")	
				applyPotion(currentPlayer)
				roundFinish=false
			} else {
				applyPotion(currentPlayer)
				printMessage("[Joueur "+string(currentPlayer.playerID)+"] Utilise une potion ! "+currentPlayer.mainPokemon+""+" +2 HP]")	
				roundFinish=true
			}
		}

		//Changer de pokémon
		if choice=="change"{
			pokemonSelected:=selectPokemon(currentPlayer)
			if pokemonSelected!="back"{
				if pokemonSelected==currentPlayer.mainPokemon{
   					printMessage("[Joueur "+string(currentPlayer.playerID)+"] "+pokemonSelected+" est déjà en combat!")
					continue
				}
				if !isPokemonAlive(pokemonSelected, currentPlayer){
					printMessage("[Joueur "+string(currentPlayer.playerID)+"] "+pokemonSelected+" n'est plus en état de se battre!")
					continue
				}
				currentPlayer.mainPokemon=pokemonSelected
				roundFinish=true
				printMessage("[Joueur "+string(currentPlayer.playerID)+"] "+pokemonSelected+" Go!")
			} 
		}
		
		//Quitter le jeu
		if choice=="quit"{
			os.Exit(0)
		}


		if roundFinish{
			if currentPlayer.playerID=='1'{
				currentPlayer=&player2
			} else {
				currentPlayer=&player1
			}
		}
	}
}
