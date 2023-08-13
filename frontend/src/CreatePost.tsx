import { Show, createSignal } from "solid-js"
import { CreateNewPost } from "../wailsjs/go/main/App"
import { reRender } from "./PostFeed"

export function CreatePost() {
    const [content, setContent] = createSignal("")
    const [flash, setFlash] = createSignal("")

    const createFn = () => {
        CreateNewPost(content())
            .then(() => { setFlash("Success!") })
            .then(reRender)
            .catch((err) => { setFlash(err) })
            .finally(() => {setInterval(() => setFlash(""),3000)})
    }

    return (
        <div class="flex flex-col">
            <input class="editor block h-32 m-2 p-4 border rounded-lg sm:text-md placeholder:text-xl bg-transparent border-gray-600 placeholder-gray-400 text-white focus:ring-blue-500 focus:border-blue-500"
                onClick={(e) => { setContent("") }}
                onInput={(e) => { setContent(e.target.value) }}
                onkeypress={(e) => {
                    if (e.key == "Enter") {
                        setContent(content().trim())
                        createFn()
                        setContent("")
                    }
                }}
                placeholder="What's happening?!"
                value={content()}
            >
            </input>
            <div class="self-end mx-2 my-2 bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-full ">
                <button onClick={createFn}>Post</button>
            </div>

            <Show when={flash() != ""}>
                <div class="text-md text-red-400">
                    {flash()}
                </div>
            </Show>
        </div>
    )
}