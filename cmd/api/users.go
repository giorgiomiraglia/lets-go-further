package main

import (
	"errors"
	"greenlight/internal/data"
	"greenlight/internal/validator"
	"net/http"
	"time"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	u := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err = u.Password.SetHash(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateUser(v, u); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(u)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicatedEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	t, err := app.models.Tokens.New(u.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	data := map[string]any{
		"userID":          u.ID,
		"activationToken": t.PlainText,
	}

	// Goroutine to send email in the background with panic recovering
	app.background(func() {
		err := app.mailer.Send(u.Email, "user_welcome.tmpl", data)
		if err != nil {
			app.logger.PrintError(err, nil)
		}
	})

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": u}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {

}
