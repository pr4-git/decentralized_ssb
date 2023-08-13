import type { Component } from 'solid-js';
import { PostFeed } from './PostFeed';
import { CreatePost } from './CreatePost';


export const MainView: Component = () => {
  return (
    <div class="app-container w-[100%] h-max grid grid-cols-7 grid-rows-8 gap-1">
      <div class='spacer col-start-2 col-span-5 row-start-1 h-3'></div>
      <div class='border col-start-2 col-span-5 row-start-2 m-2'>
      <CreatePost/>
      </div>
      <div class='row-start-3 col-start-2 col-span-5 m-2'>
      <PostFeed/>
      </div>
    </div>
  );
};

