package ssp

import (
	"errors"
	"fmt"
	"net/http"

	sqrl "github.com/RaniSputnik/sqrl-go"
)

const xFormURLEncoded = "application/x-www-form-urlencoded"

var v1Only = []string{sqrl.V1}

func clientFailure(response *sqrl.ServerMsg) {
	response.Set(sqrl.TIFCommandFailed).Set(sqrl.TIFClientFailure)
}

func serverError(response *sqrl.ServerMsg) {
	response.Set(sqrl.TIFCommandFailed).Set(sqrl.TIFTransientError)
}

// TODO: This method is ridiculously large, we should be able to break it down
// and move some of the functionality (particularly validation) to the core SQRL
// package for folks who don't need a SSP server.
func (server *Server) ClientHandler(store Store, tokens TokenGenerator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		server.logger.Printf("Got SQRL request: %v\n", r)

		response := genNextResponse(server, r)
		defer writeResponse(w, response)

		req, err := parseRequest(r)
		if err != nil {
			server.logger.Printf("Client failure, %v'\n", err)
			clientFailure(response)
			return
		}

		firstTransaction, err := store.GetFirstTransaction(ctx, req.Nut)
		if err != nil {
			server.logger.Printf("Failed to retrieve first transaction: %v\n", err)
			serverError(response)
			return
		}

		// This is not ideal positioning for performance. Ideally
		// we would _only_ GetFirstTransaction if the basic
		// verify checks succeeded (ie. signatures valid etc.)
		// but that would require breaking Verify up into multiple
		// steps. Probably ok as the transaction store is likely to
		// entirely stored in cache.
		client, err := sqrl.Verify(req, firstTransaction, response)
		if err != nil {
			server.logger.Printf("Failed to verify transaction: %v", err)
			return
		}

		if err := store.SaveTransaction(ctx, &sqrl.Transaction{
			Request: req,
			Next:    response.Nut,
		}); err != nil {
			server.logger.Printf("Failed to save transaction: %v\n", err)
			serverError(response)
			return
		}

		// TODO: Pass previous identities to "GetByIdentity"
		currentUser, err := store.GetUserByIdentity(ctx, client.Idk)
		if err != nil {
			server.logger.Printf("Failed to determine if identity is known: %v\n", err)
			serverError(response)
			return
		} else if currentUser != nil {
			response.Set(sqrl.TIFCurrentIDMatch)
		}

		switch client.Cmd {
		case sqrl.CmdIdent:
			// Create user if they do not already exist
			if currentUser == nil {
				currentUser, err = store.CreateUser(ctx, client.Idk)
				if err != nil {
					server.logger.Printf("Failed to create user: %v\n", err)
					serverError(response)
					return
				}
			}

			// Generate a new token that can be exchanged for user credentials
			// TODO: It would be great if we could guarantee the size of tokens
			// for DB backends that want to specify the column size for the token
			token := tokens.Token(currentUser.Id)
			// Record that this transaction was a success, store the token
			sessionID := req.Nut
			if firstTransaction != nil {
				sessionID = firstTransaction.Nut
			}
			err = store.SaveIdentSuccess(r.Context(), sessionID, token)
			if err != nil {
				server.logger.Printf("Failed to save ident success: %v\n", err)
				serverError(response)
				return
			}

			if client.HasOpt(sqrl.OptCPS) {
				response.URL = getTokenRedirectURL(server, token)
			}
		case sqrl.CmdQuery:
			// TODO: Anything need to be done here?

		default:
			// In all other cases, not supported
			response.Set(sqrl.TIFFunctionNotSupported)
		}
	})
}

func genNextResponse(server *Server, r *http.Request) *sqrl.ServerMsg {
	nextNut := server.Nut()
	return &sqrl.ServerMsg{
		Ver: v1Only,
		Nut: nextNut,
		Qry: server.clientEndpoint + "?nut=" + string(nextNut),
	}
}

func parseRequest(r *http.Request) (*sqrl.Request, error) {
	if contentType := r.Header.Get("Content-Type"); contentType != xFormURLEncoded {
		return nil, fmt.Errorf("invalid content type: '%s'", contentType)
	}

	nut := sqrl.Nut(r.URL.Query().Get("nut"))
	if nut == "" {
		return nil, errors.New("missing required parameter: 'nut'")
	}

	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	clientRaw := r.Form.Get("client")
	serverRaw := r.Form.Get("server")
	ids := sqrl.Signature(r.Form.Get("ids"))
	// TODO: pids

	return &sqrl.Request{
		Nut:      nut,
		Client:   clientRaw,
		Server:   serverRaw,
		Ids:      ids,
		ClientIP: ClientIP(r),
	}, nil
}

func writeResponse(w http.ResponseWriter, response *sqrl.ServerMsg) {
	encoded, err := response.Encode()
	if err != nil {
		panic(err)
	}
	// TODO: This is a bit janky but it's what the reference
	// implementation does. Should probably question the use
	// of this content type given it's not in the form key=value.
	w.Header().Set("Content-Type", xFormURLEncoded)
	if _, err := w.Write([]byte(encoded)); err != nil {
		panic(err) // TODO: What to do here?
	}
}
