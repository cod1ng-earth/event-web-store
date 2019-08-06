module Main exposing (main)

-- https://package.elm-lang.org/packages/NoRedInk/elm-sweet-poll/latest/
-- https://package.elm-lang.org/packages/simonh1000/elm-jwt/latest/
-- https://github.com/simonh1000/elm-jwt/blob/master/examples/node/src/Main.elm

import Browser
import Browser.Events exposing (onVisibilityChange, Visibility(..))
import Html exposing (Html, button, div, text, h1, ul, li, a)
import Html.Events exposing (onClick)
import Html.Attributes exposing (href)
import Time
import Task
import Random
import Http
import Delay
import Json.Decode exposing (Decoder, map3, float, string, field, list, decodeString)


type alias Model =
    { count : Int
    , same : Int
    , text : String
    , out : String
    , error : String
    , products : Products
    }


init : () -> ( Model, Cmd Msg )
init _ =
    ( { count = 0
      , same = 0
      , text = ""
      , out = ""
      , error = ""
      , products = []
      }
    , Cmd.batch [ fetchResults ]
    )


type Msg
    = Increment
    | Decrement
    | LoadProducts
    | Tick Time.Posix
    | Reload
    | NewFace Int
    | GotText (Result Http.Error String)
    | GotProducts (Result Http.Error Products)
    | VisibilityChange Browser.Events.Visibility


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        Tick _ ->
            ( model, Random.generate NewFace roll )

        Reload ->
            ( model, fetchResults )

        LoadProducts ->
            ( model, fetchProducts )

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
                    ( { model | error = toString e }, Delay.after 500 Delay.Millisecond Reload )

        GotProducts result ->
            case result of
                Ok fullText ->
                    ( { model | products = fullText }, Cmd.none )

                Err e ->
                    ( { model | error = toString e }, Cmd.none )
        
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
            [ button [ onClick LoadProducts ] [ text "load products" ]
            , button [ onClick Increment ] [ text "+1" ]
            , div [] [ text <| String.fromInt model.count ]
            , div [] [ text <| String.fromInt model.same ]
            , button [ onClick Decrement ] [ text "-1" ]
            , div [] [ text model.text ]
            , div [] [ text model.out ]
            , div [] [ text model.error ]
            , div [] [ renderProducts model.products ]
            ]
        ]

renderProducts : Products -> Html msg
renderProducts lst =
    ul []
        (List.map (\l -> li [] [ renderProduct l ]) lst )

renderProduct : Product  -> Html msg
renderProduct product =
    a [ href ("/product/" ++ product.uuid) ] [text product.title]

roll : Random.Generator Int
roll =
    Random.int 1 6


fetchResults : Cmd Msg
fetchResults =
    Http.get
        { url = "http://localhost:8080/foobarbaz"
        , expect = Http.expectString GotText
        }


fetchProducts : Cmd Msg
fetchProducts =
    Http.get
        { url = "http://localhost:8080/all-products"
        , expect = Http.expectJson GotProducts productsDecoder
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

type alias Products = List Product

productsDecoder : Decoder Products
productsDecoder = list productDecoder

type alias Product =
  { uuid             : String
  , title            : String
  , price            : Float
  }

productDecoder : Decoder Product
productDecoder = map3 Product
  (field "uuid"          string)
  (field "title"         string)
  (field "price"         float)
