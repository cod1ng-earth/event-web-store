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
    { error : String
    , products : Products
    }


init : () -> ( Model, Cmd Msg )
init _ =
    ( { error = ""
      , products = []
      }
    , Cmd.batch [ fetchProducts ]
    )


type Msg
    = LoadProducts
    | GotProducts (Result Http.Error Products)


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        LoadProducts ->
            ( model, fetchProducts )

        GotProducts result ->
            case result of
                Ok fullText ->
                    ( { model | products = fullText, error = "" }, Cmd.none )

                Err e ->
                    ( { model | error = toString e }, Cmd.none )


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


fetchProducts : Cmd Msg
fetchProducts =
    Http.get
        { url = "http://localhost:8080/all-products"
        , expect = Http.expectJson GotProducts productsDecoder
        }


subscriptions : Model -> Sub Msg
subscriptions _ = Sub.batch []

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
