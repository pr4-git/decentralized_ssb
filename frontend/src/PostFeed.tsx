import { For, Show, createResource, createSignal, createMemo, onMount } from "solid-js";
import { ViewAllPosts } from "../wailsjs/go/main/App";
import { decodeTime } from 'ulid';
import { ProfileViewer, UpdateProfile } from "./ProfileViewer";
import { Select } from "@thisbeyond/solid-select";

export const [renderProfiles, setRenderProfiles] = createSignal(0)
export function reRender() {
    setRenderProfiles(renderProfiles() + 1)
}

export function PostFeed() {

    const TimeFromULID = (ulid: string) => {
        let timestamp = decodeTime(ulid)
        let date = new Date(timestamp)
        return date
    }

    const [posts, { mutate, refetch }] = createResource(renderProfiles, ViewAllPosts)

    onMount(reRender)
    setInterval(() => {
        UpdateProfile()
        reRender()
    }, 5000)

    const [ViewMode, setViewMode] = createSignal("all")
    const filterByPk = (pubkey: string) => posts()?.filter(post => post.author.toString() == pubkey)
    const authorList = createMemo(() => {
        let unique: string[] = []
        let authors = (posts()?.map(post => post.author.toString()) || [])
        authors.forEach(author => {
            if (!unique.includes(author))
                unique.push(author)
        })
        return unique
    })
    const [value, setValue] = createSignal(null)
    const [initialValue, setInitialValue] = createSignal(null, { equals: false })

    return (
        <div>
            <div class="grid grid-cols-[15%_70%_15%] grid-rows-[30%_90%] outline-none ">
                <div class="col-start-1 col-span-1 row-start-1 row-span-1">Filter</div>
                <div class="col-start-2 col-span-1 row-start-1">
                    <Select
                        class="selection text-xs border bg-transparent "
                        initialValue={initialValue()}
                        options={authorList()}
                        onChange={setValue}
                    />
                </div>
                <div class="col-start-3 col-span-1 row-start-1 row-span-1 flex place-content-around">
                    <button onClick={() => setInitialValue(null)}>❌</button>
                </div>

            </div>
            <div class="feed-view h-max flex flex-col-reverse gap-4">
                <For each={value() ? filterByPk(value()) : posts()}>{post =>
                    <div class="grid grid-cols-[23%_77%] grid-rows-[70%_30%] gap-2 border h-48 my-2  p-3">
                        <div class="profile-icon col-start-1 row-start-1 col-span-1 row-span-2 flex flex-col justify-center items-center">
                            <ProfileViewer pubkey={post.author.toString()} />
                        </div>
                        <div class="post-content borpx-3 text-xl col-start-2 row-start-1 col-span-1 row-span-1 flex items-center">
                            {post.content}
                        </div>
                        <div class="px-2 row-start-2 col-start-2 row-span-1 col-span-1 flex flex-col items-end justify-center">
                            <div class="post-id text-xs">
                                hash: {post.hash.slice(0, 16)}
                                <br />
                                parent: {post.parent?.slice(0, 16) || "initial post"}
                            </div>
                            <div class="post-time text-xs">
                                {TimeFromULID(post.id).toLocaleString()}
                            </div>
                        </div>
                    </div>
                }
                </For>
                <Show when={posts.error}>
                    <div>Error loading posts. {posts.error}</div>
                </Show>
            </div>
        </div>
    )

}

