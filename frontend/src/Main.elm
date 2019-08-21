module Main exposing (main)

-- https://package.elm-lang.org/packages/NoRedInk/elm-sweet-poll/latest/
-- https://package.elm-lang.org/packages/simonh1000/elm-jwt/latest/
-- https://github.com/simonh1000/elm-jwt/blob/master/examples/node/src/Main.elm
-- https://github.com/rtfeldman/elm-spa-example.git
-- https://medium.com/@_rchaves_/child-parent-communication-in-elm-outmsg-vs-translator-vs-nomap-patterns-f51b2a25ecb1
-- https://github.com/evancz/elm-todomvc.git
-- https://elmprogramming.com/saving-app-state.html
-- https://package.elm-lang.org/packages/elm/browser/latest/Browser-Navigation
-- https://package.elm-lang.org/packages/billstclair/elm-websocket-client/latest/
-- https://riptutorial.com/elm/example/13036/program-with-flags
-- http://blog.jenkster.com/2016/04/how-i-structure-elm-apps.html

import Browser
import Message exposing (Msg)
import Model exposing (init, Model)
import View exposing (view)
import Update exposing (update)


subscriptions : Model -> Sub Msg
subscriptions _ =
    Sub.batch []


main : Program () Model Msg
main =
    Browser.element
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        }
