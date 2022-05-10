package api

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

type Item struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Server struct {
	*mux.Router
	shoppingItems []Item
}

func NewServer() *Server {
	// C'est exactement le même délire que $s = new Server en PHP
	// Du coup, le retour est du type *Server, d'où la suite
	s := &Server{
		Router:        mux.NewRouter(),
		shoppingItems: []Item{},
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.HandleFunc("/items", s.listShoppingItem()).Methods("GET")
	s.HandleFunc("/items", s.createShoppingItem()).Methods("POST")
	s.HandleFunc("/items/{id}", s.removeShoppingItem()).Methods("DELETE")
}

// Une fonction qui retourne une closure
func (s *Server) createShoppingItem() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var i Item

		/* Je crée un nouveau décodeur JSON qui va lire le body de ma requête HTTP
		ensuite je lui demande de mettre le résultat du décodage dans ma variable i fraichement crée
		d'où le passage par référence
		Le reste c'est de la gestion d'erreur, si j'ai une erreur je renvoie une réponse HTTP d'erreur avec
		le message et un code de réponse 400
		*/
		if err := json.NewDecoder(request.Body).Decode(&i); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		i.ID = uuid.New()
		s.shoppingItems = append(s.shoppingItems, i)

		writer.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(writer).Encode(i); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

func (s *Server) listShoppingItem() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(writer).Encode(s.shoppingItems); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

func (s *Server) removeShoppingItem() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		idStr, _ := mux.Vars(request)["id"]
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		for i, item := range s.shoppingItems {
			if item.ID == id {
				s.shoppingItems = append(s.shoppingItems[:i], s.shoppingItems[i+1:]...)
				break
			}
		}
	}
}
