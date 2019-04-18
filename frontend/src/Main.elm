module Main exposing (main)

-- https://package.elm-lang.org/packages/NoRedInk/elm-sweet-poll/latest/
-- https://package.elm-lang.org/packages/simonh1000/elm-jwt/latest/
-- https://github.com/simonh1000/elm-jwt/blob/master/examples/node/src/Main.elm

import Browser
import Browser.Events exposing (onVisibilityChange, Visibility(..))
import Html exposing (Html, button, div, text, h1)
import Html.Events exposing (onClick)
import Time
import Task
import Random
import Http
import Delay


type alias Model =
    { count : Int
    , same : Int
    , text : String
    , out : String
    }


init : () -> ( Model, Cmd Msg )
init _ =
    ( { count = 0
      , same = 0
      , text = ""
      , out = ""
      }
    , Cmd.batch [ fetchResults ]
    )


type Msg
    = Increment
    | Decrement
    | Tick Time.Posix
    | Reload
    | NewFace Int
    | GotText (Result Http.Error String)
    | VisibilityChange Browser.Events.Visibility


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        Tick _ ->
            ( model, Random.generate NewFace roll )

        Reload ->
            ( model, fetchResults )

        NewFace role ->
            if model.count == role then
                ( { model | same = model.same + 1 }, Cmd.none )
            else
                ( { model | count = role }, Cmd.none )

        Increment ->
            ( { model | count = model.count + 1 }, Cmd.none )

        Decrement ->
            if model.count > 0 then
                ( { model | count = model.count - 1 }, Cmd.none )
            else
                ( model, Cmd.none )

        GotText result ->
            case result of
                Ok fullText ->
                    ( { model | text = fullText }, Delay.after 500 Delay.Millisecond Reload )

                Err e ->
                    ( { model | text = toString e }, Delay.after 500 Delay.Millisecond Reload )
        
        VisibilityChange Visible ->
            ( { model | out = "v" }, Cmd.none )

        VisibilityChange Hidden ->
            ( { model | out = "h" }, Cmd.none )



toString : Http.Error -> String
toString error =
    case error of
        Http.BadUrl url ->
            url ++ " is bad"

        Http.Timeout ->
            "Timeout"

        Http.NetworkError ->
            "Network Error"

        Http.BadStatus status ->
            String.fromInt status ++ " Status"

        Http.BadBody msg ->
            msg


view : Model -> Html Msg
view model =
    div []
        [ h1 [] [ text "Hello, world." ]
        , div []
            [ button [ onClick Increment ] [ text "+1" ]
            , div [] [ text <| String.fromInt model.count ]
            , div [] [ text <| String.fromInt model.same ]
            , button [ onClick Decrement ] [ text "-1" ]
            , div [] [ text model.text ]
            , div [] [ text model.out ]
            ]
        ]


roll : Random.Generator Int
roll =
    Random.int 1 6


fetchResults : Cmd Msg
fetchResults =
    Http.get
        { url = "http://localhost:8080/foobarbaz"
        , expect = Http.expectString GotText
        }


subscriptions : Model -> Sub Msg
subscriptions _ =
    Sub.batch
        [ Time.every 1000 Tick
        -- , Time.every 2000 Reload
        , onVisibilityChange abc
        ]

abc : Browser.Events.Visibility -> Msg
abc visibility = VisibilityChange visibility

main : Program () Model Msg
main =
    Browser.element
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        }
