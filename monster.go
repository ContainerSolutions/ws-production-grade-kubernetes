package main

import (
    "fmt"
    "image"
    "image/color"
    "math/rand"
)

func FetchMonster(id string, size int) image.Image {

    nameBytes := []byte(id)
    avatar := image.NewRGBA(image.Rect(0, 0, size, size))
    bgColor := CalcBGColor(nameBytes)
    PaintBG(avatar, bgColor, size)
    PaintMonster(avatar, nameBytes, CalcPixelColor(nameBytes), bgColor, size)
    fmt.Println("fetch monster")
    return avatar
}

func PaintMonster(avatar *image.RGBA, nameBytes []byte, pixelColor color.RGBA, bgColor color.RGBA, size int) {
    
    var nameSum int64
    for i := range nameBytes {
        nameSum += int64(nameBytes[i])
    }

    rand.Seed(nameSum)

    // Avatar random parts.
    var parts = []MonsterPart {
        Monsters.legs[rand.Intn(len(Monsters.legs))],
        Monsters.hair[rand.Intn(len(Monsters.hair))],
        Monsters.arms[rand.Intn(len(Monsters.arms))],
        Monsters.body[rand.Intn(len(Monsters.body))],
        Monsters.eyes[rand.Intn(len(Monsters.eyes))],
        Monsters.mouth[rand.Intn(len(Monsters.mouth))],
    }

    // Fill avatar with random parts.
    for _, i := range parts {
        drawPart(i, avatar, size, pixelColor, bgColor);
    }
}

func drawPart(part MonsterPart, avatar *image.RGBA, size int, pixelColor color.RGBA, bgColor color.RGBA) {

    dotSize := (size / len(part.piece))

    for r, row := range part.piece {
        y := r * dotSize
        for c, col := range row.part {
            x := c * dotSize
            switch col {
            case 1:
                drawRect(avatar, dotSize, x, y, pixelColor)
            case 2:
                drawRect(avatar, dotSize, x, y, bgColor)
            }
        }
    }
}

func drawRect(avatar *image.RGBA, size int, x int, y int, pixelColor color.RGBA) {
    for i := x; i < x + size; i++ {
        for j := y; j < y + size; j++ {
            avatar.SetRGBA(i, j, pixelColor)
        }
    }
}

func CalcPixelColor(nameBytes []byte) (pixelColor color.RGBA) {
    pixelColor.A = 255

    var mutator = byte((len(nameBytes) * 4))

    pixelColor.R = nameBytes[0] * mutator
    pixelColor.G = nameBytes[1] * mutator
    pixelColor.B = nameBytes[2] * mutator

    return
}

func CalcBGColor(nameBytes []byte) (bgColor color.RGBA) {
    bgColor.A = 255

    var mutator = byte((len(nameBytes) * 2))

    bgColor.R = nameBytes[0] * mutator
    bgColor.G = nameBytes[1] * mutator
    bgColor.B = nameBytes[2] * mutator

    return
}

func PaintBG(avatar *image.RGBA, bgColor color.RGBA, size int) {
    for y := 0; y < size; y++ {
        for x := 0; x < size; x++ {
            avatar.SetRGBA(x, y, bgColor)
        }
    }
}