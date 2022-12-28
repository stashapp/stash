export function hexToRgb (hex:string) {
    if(hex == null) {
        return null
    }
    
    const split = hex.replace(/^#?([a-f\d])([a-f\d])([a-f\d])$/i, (m, r, g, b) => '#' + r + r + g + g + b + b).substring(1).match(/.{2}/g)
    if(split != null && split.length > 0)
    {
        const ret = split.map(x => parseInt(x, 16))
        return ret
    }
    return null
}

export function contrastingTextColor(backgroundColorHex:string){
    const backgroundColorRGB = hexToRgb(backgroundColorHex)
    if(backgroundColorRGB == null){
        return "black"
    }

    return ((backgroundColorRGB[0]*299)+(backgroundColorRGB[1]*587)+(backgroundColorRGB[2]*114))/1000 >= 128 ? "#182026" : "#f5f8fa"

}