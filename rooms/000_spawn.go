package rooms

import "fmt"

func init() {
	Rooms.room["000_spawn"] = Room{
		tiles: `// First line is ignored because lazy.
		##D##
		# . #
		#...C
		# . #
		#####
		`,
		tileMap: map[string]string{
			"#": "wall",
			".": "floor",
			"D": "door",
			"C": "computer",
		},
		entities: `
		     
		     
		  @  
		     
		     
		`,
		entityMap: map[string]string{
			"@": "player",
		},
		metadata: make(map[string]interface{}),
		enter: func() {
			fmt.Println("entered spawn")
		},
		leave: func() {
			fmt.Println("left spawn")
		},
		update: func() {
			//
		},
	}
}
