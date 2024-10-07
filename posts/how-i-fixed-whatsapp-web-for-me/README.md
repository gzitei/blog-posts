# How I fixed WhatsApp Web (for me)

### The Problem

To fit all my developer tools and documentations in my screen, I like to use a tilling window tool called [Tactile](https://gitlab.com/lundal/tactile) on Ubuntu. That's where the problem starts.

Even though I'm pretty sure people at Meta must have a good reason for that, but WhatsApp Web's responsiveness leaves room for improvement.

I like to have it tilled as a little vertical bar side by side with my text editor so I can answer my girlfriend even when I am working (sorry my friends, she takes all my WhatsApp attention).

That's where responsiveness falls short. The contact's sidebar takes almost all the viewport, there is no space left for the user to interact with the conversation, as shown in Figure 1.

![Figure 1: WhatsApp Web responsiveness issue](https://raw.githubusercontent.com/gzitei/blog-posts/refs/heads/main/posts/how-i-fixed-whatsapp-web-for-me/images/problem.webp)

### The Solution

I tried some chrome extensions and even though they solved the sidebar problem, some of them left me with another problem: I had to turn off the extension in order to access my contact list. That got me so frustrated, I had to solve this problem myself, but I had never built a browser extension before...

That's when I started looking for some github repositories and I came across [this repo](https://github.com/julianilevy/whatsapp-web-sidebar-toggle) which I decided to fork and iterate over the solution.

I needed the contact bar to be hidden, using devtools I started exploring the HTML tree and I found that I needed to hide the parent of an element which id was `side`: 

```javascript
const sidebar = document.querySelector("#side")?.parentElement;
if (sidebar) sidebar.style.display = "none";
```

This solved the problem of the sidebar, but I still would want to be able to choose another contact from my list and also use the other WhatsApp functionalities, so I had to place event listeners to the page in order to toggle the sidebar visibility.

First, I defined a `toggleSideBar` function, which is responsbible for turning on and off the sidebar visibility:

```javascript
const toggleSideBar = () => {
  if (sidebar) {
    sidebar.style.display =
      sidebar.style.display === "none" ? "" : "none";
  }
};
```
When I started experimenting with `toggleSideBar`, I noticed there are two different elements used as sidebar: one for the contact list and another for the other functionalities like settings, community, status and newsletters. So I had to define another sidebar and we all know: naming things is really hard, so I named this element `secondarybar` and a `toggleSecondarySideBar` function:

```javascript
const secondarybar = document.querySelector("[class='`aohf']");
if (secondarybar) secondarybar.style.display = "none";

const toggleSecondarySideBar = () => {
  if (secondarybar) {
    if (sidebar.style.display != "none") {
      sidebar.style.display = "none";
    }
    secondarybar.style.display =
      secondarybar.style.display === "none" ? "" : "none";
  }
};
```

Also, I had to update `toggleSideBar`:

```javascript
const toggleSideBar = () => {
  if (sidebar) {
    if (secondarybar.style.display != "none") {
      secondarybar.style.display = "none";
    }
    sidebar.style.display =
      sidebar.style.display === "none" ? "" : "none";
  }
};
```

Then, using the DOM's concept of `event bubbling`, widely used in React, I placed the event listener in the document's body, using the event's target to define whether to toggle the sidebar visibility or not:

```javascript
document.body.addEventListener("click", (e) => {
  let label = e.target?.dataset?.icon;
  if (!label) return;

  let selected = label.split("-")[0];
  if (!views.includes(selected)) return;

  if (currentView == "") {
    if (selected == views[0]) {
      toggleSideBar();
    } else {
      toggleSecondarySideBar();
    }
    currentView = selected;
    return;
  }

  if (selected == currentView) {
    if (selected == views[0]) {
      toggleSideBar();
    } else {
      toggleSecondarySideBar();
    }
    currentView = "";
    return;
  }

  if (currentView == views[0] && selected != views[0]) {
    toggleSideBar();
    toggleSecondarySideBar();
    currentView = selected;
    return;
  }

  if (currentView != views[0] && selected == views[0]) {
    toggleSecondarySideBar();
    toggleSideBar();
    currentView = selected;
    return;
  }

  currentView = selected;
});
```
With the current implementation, I have all I need! Now, I am able to use WhatsApp side by side with my editor, I just need to install my extension and... OH NO! The text input field looks awful!

![Figure 2: Input field is broken](https://raw.githubusercontent.com/gzitei/blog-posts/refs/heads/main/posts/how-i-fixed-whatsapp-web-for-me/images/input-bar.webp)

So I decided to fix the chat window element width:

```javascript
const main = document.querySelectorAll("[class='_aigv _aigz']")[1];
if (main) {
  main.style.maxWidth = `calc(100vw - var(--navbar-width))`;
}
```
By setting the maximum width of the chat window element to fill the available viewport discounted by the lateral navbar width, the text input field was fixed!

![Figure 3: Input field is broken](https://raw.githubusercontent.com/gzitei/blog-posts/refs/heads/main/posts/how-i-fixed-whatsapp-web-for-me/images/fixed-input-bar.webp)

Now I have everything working perfectly to fit my workflow! The next image presents the final result with the extension running.

![Figure 4: Final result](https://raw.githubusercontent.com/gzitei/blog-posts/refs/heads/main/posts/how-i-fixed-whatsapp-web-for-me/images/final-result.webp)

I finished by implementing some other adjustments to improve on the apps UI in fullscreen mode, removing some unnecessary margins around the main content.

If you are interested, you may check my [github repository](https://github.com/gzitei/whatsapp-web-sidebar-toggle/), there are some instructions on how to use the extension.

Also, if you have suggestions or feedback about this project, feel free to reach out to me on my social media.
