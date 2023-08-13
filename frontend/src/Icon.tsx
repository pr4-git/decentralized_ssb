
type IconProps = { iconSrc: string, msg: string }

export function Icon({ iconSrc, msg }: IconProps) {
    return (
        <div class="flex flex-col justify-around">
            <div class="flex flex-row place-content-center">
                <div class="border rounded-full p-3">
                    <img src={iconSrc} color="white" width="40" height="40"></img>
                </div>
            </div>
            <div class="m-2 text-center">{msg}</div>
        </div>
    )
}