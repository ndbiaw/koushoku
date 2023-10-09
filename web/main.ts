import {
  getReaderSettings,
  getSettings,
  Mode,
  readerSettings,
  setReaderSettings,
  setSettings,
  settings
} from "./settings";
import "./styles/main.less";

interface PageState {
  isViewing?: boolean;
  isPreloading?: boolean;
  isPreloaded?: boolean;
}

const pageStates: PageState[] = [];
const maxPreloads = 3;

let id: string;
let slug: string;
let origin: string;
let pageContainer: HTMLAnchorElement;

let totalPages: number;
let currPageSpans: NodeListOf<HTMLElement>;
let currPageNum: number;
let currPageImage: HTMLImageElement;

let firstPageAnchors: NodeListOf<HTMLAnchorElement>;
let lastPageAnchors: NodeListOf<HTMLAnchorElement>;
let prevPageAnchors: NodeListOf<HTMLAnchorElement>;
let nextPageAnchors: NodeListOf<HTMLAnchorElement>;

const mutex = { current: false };
const preventDefault = (ev: MouseEvent) => {
  ev.preventDefault();
  ev.stopPropagation();
  ev.stopImmediatePropagation();
};

const applyReaderSettings = () => {
  if (readerSettings.maxWidth > 0) {
    pageContainer.style.maxWidth = `${readerSettings.maxWidth}px`;
  } else pageContainer.style.maxWidth = "";

  if (readerSettings.zoomLevel > 0) {
    (pageContainer.style as any).zoom = `${readerSettings.zoomLevel}`;
  } else (pageContainer.style as any).zoom = "";

  const img = pageContainer.firstElementChild as HTMLImageElement;
  if (readerSettings.zoomLevel > 1.0) {
    img.style.maxWidth = `${window.innerWidth * readerSettings.zoomLevel}px`;
  } else img.style.maxWidth = "";

  document.querySelectorAll(".zoom-level").forEach((span: HTMLElement) => {
    span.innerText = `${(readerSettings.zoomLevel * 100).toFixed()}%`;
  });
};

let showSettingsPopup: () => void;
const attachHandlers = () => {
  const zoomInBtns = document.querySelectorAll(".zoom-in") as NodeListOf<HTMLButtonElement>;
  const zoomOutBtns = document.querySelectorAll(".zoom-out") as NodeListOf<HTMLButtonElement>;
  const toggleBtns = document.querySelectorAll(".toggle-settings") as NodeListOf<HTMLButtonElement>;

  zoomInBtns.forEach(btn => {
    if (btn.dataset.attached) return;
    btn.dataset.attached = "true";
    btn.addEventListener("click", () => {
      const zoomLevel = (readerSettings.zoomLevel + 0.1).toFixed(1);

      setReaderSettings("zoomLevel", Math.max(0.1, Number(zoomLevel)));
      applyReaderSettings();
    });
  });

  zoomOutBtns.forEach(btn => {
    if (btn.dataset.attached) return;
    btn.dataset.attached = "true";
    btn.addEventListener("click", () => {
      const zoomLevel = (readerSettings.zoomLevel - 0.1).toFixed(1);
      setReaderSettings("zoomLevel", Math.min(2.0, Number(zoomLevel)));
      applyReaderSettings();
    });
  });

  toggleBtns.forEach(btn => {
    if (btn.dataset.attached) return;
    btn.dataset.attached = "true";
    btn.addEventListener("click", () => showSettingsPopup());
  });
};

showSettingsPopup = () => {
  if (document.getElementById("settings-popup")) {
    return;
  }

  const popup = document.createElement("div");
  popup.id = "settings-popup";

  const overlay = document.createElement("div");
  overlay.addEventListener("click", () => popup.remove());
  overlay.classList.add("settings-overlay");

  const content = document.createElement("div");
  content.classList.add("settings-content");

  const header = document.createElement("div");
  header.classList.add("settings-header");

  const title = document.createElement("span");
  title.textContent = "Settings";
  header.appendChild(title);

  const container = document.createElement("div");
  container.classList.add("settings-container");
  {
    const section = document.createElement("div");
    section.classList.add("settings-section");

    const label = document.createElement("span");
    label.textContent = "Zoom";

    const actions = document.createElement("div");

    const zoomOutBtn = document.createElement("button");
    zoomOutBtn.type = "button";
    zoomOutBtn.textContent = "-";
    zoomOutBtn.classList.add("zoom-out");

    const zoomLevelText = document.createElement("span");
    zoomLevelText.textContent = `${(readerSettings.zoomLevel * 100).toFixed()}%`;
    zoomLevelText.classList.add("zoom-level");

    const zoomInBtn = document.createElement("button");
    zoomInBtn.type = "button";
    zoomInBtn.textContent = "+";
    zoomInBtn.classList.add("zoom-in");

    actions.append(zoomOutBtn, zoomLevelText, zoomInBtn);
    section.append(label, actions);
    container.appendChild(section);
  }
  {
    const section = document.createElement("div");
    section.classList.add("settings-section", "max-width");

    const label = document.createElement("span");
    label.textContent = "Max. width";

    const actions = document.createElement("div");
    const wrapper = document.createElement("div");
    wrapper.classList.add("wrapper");

    const input = document.createElement("input");
    input.type = "number";
    input.min = "0";
    input.defaultValue = readerSettings.maxWidth.toString();
    input.addEventListener("change", () => {
      const value = Number(input.value);
      if (Number.isNaN(value)) {
        input.value = "0";
      }

      if (readerSettings.maxWidth !== value) {
        setReaderSettings("maxWidth", value);
        applyReaderSettings();
      }
    });

    const suffix = document.createElement("span");
    suffix.textContent = "px";

    wrapper.append(input, suffix);
    actions.append(wrapper);
    section.append(label, actions);
    container.appendChild(section);
  }
  /*   {
    const section = document.createElement("div");
    section.classList.add("settings-section", "reading-mode");

    const label = document.createElement("span");
    label.textContent = "Reading mode";

    const actions = document.createElement("div");
    const wrapper = document.createElement("div");
    wrapper.classList.add("wrapper");

    const normalBtn = document.createElement("button");
    normalBtn.classList.add("toggle-normal");
    normalBtn.textContent = "Normal";
    normalBtn.dataset.value = Mode.Normal.toString();
    if (settings.mode === Mode.Normal) {
      normalBtn.classList.add("active");
    }

    const stripBtn = document.createElement("button");
    stripBtn.classList.add("toggle-strip");
    stripBtn.textContent = "Strip";
    stripBtn.dataset.value = Mode.Strip.toString();
    if (settings.mode === Mode.Strip) {
      stripBtn.classList.add("active");
    }

    const onClick = (ev: MouseEvent) => {
      const mode = Number((ev.target as HTMLButtonElement).dataset.value) as Mode;
      if (settings.mode === mode) {
        return;
      }

      normalBtn.classList.toggle("active");
      stripBtn.classList.toggle("active");
      setSetting("mode", mode);
      applySettings();
    };

    normalBtn.addEventListener("click", onClick);
    stripBtn.addEventListener("click", onClick);

    wrapper.append(normalBtn, stripBtn);
    actions.appendChild(wrapper);
    section.append(label, actions);
    container.appendChild(section);
  } */

  const footer = document.createElement("div");
  footer.classList.add("footer");

  const closeBtn = document.createElement("button");
  closeBtn.classList.add("close");
  closeBtn.type = "button";
  closeBtn.textContent = "Close";
  closeBtn.addEventListener("click", () => popup.remove());

  footer.appendChild(closeBtn);
  container.appendChild(footer);
  content.append(header, container);
  popup.append(overlay, content);

  document.body.appendChild(popup);
  attachHandlers();
};

const changePage = (targetPageNum: number) => {
  if (mutex.current) return;
  mutex.current = true;

  try {
    currPageNum = targetPageNum;

    pageStates.find(p => p.isViewing).isViewing = false;
    pageStates[currPageNum - 1].isViewing = true;

    if (readerSettings.mode === Mode.Normal) {
      pageContainer.href = `/archive/${id}/${slug}/${currPageNum}`;
      const newImg = document.createElement("img");
      newImg.src = `${origin}/data/${id}/${currPageNum}.jpg`;

      currPageImage.replaceWith(newImg);
      currPageImage = newImg;
    } else {
      //
    }

    currPageSpans.forEach(span => (span.textContent = currPageNum.toString()));
    window.history.replaceState(null, "", `/archive/${id}/${slug}/${currPageNum}`);
    prevPageAnchors.forEach(e => (e.href = `/archive/${id}/${slug}/${currPageNum - 1}`));
    nextPageAnchors.forEach(e => (e.href = `/archive/${id}/${slug}/${currPageNum + 1}`));
  } finally {
    mutex.current = false;
  }
};

const initPreload = () => {
  const currentIndex = pageStates.findIndex(e => e.isViewing);

  let count = 0;
  for (let i = 0; i < pageStates.length; i++) {
    const page = pageStates[i];
    if (
      page.isPreloading ||
      page.isPreloaded ||
      i < Math.max(0, currentIndex - maxPreloads) ||
      i > Math.min(pageStates.length, currentIndex + maxPreloads)
    ) {
      continue;
    }

    const pageNum = i + 1;
    page.isPreloading = true;

    // eslint-disable-next-line @typescript-eslint/no-loop-func
    setTimeout(() => {
      const img = document.createElement("img");
      img.src = `${origin}/data/${id}/${pageNum}.jpg`;

      const onComplete = () => {
        page.isPreloading = false;
        page.isPreloaded = true;
      };
      const onFailed = () => (page.isPreloading = false);

      if (img.complete) {
        onComplete();
      } else {
        img.addEventListener("load", onComplete);
        img.addEventListener("error", onFailed);
      }
    }, 1000 * count);
    count++;
  }
};

const applySettings = () => {
  document.querySelectorAll("article[data-hidden]").forEach(e => e.removeAttribute("data-hidden"));
  if (settings.blacklist) {
    document.querySelectorAll("article[data-tags]").forEach((e: HTMLElement) => {
      const tags = (e.dataset.tags as string).trim().toLowerCase();
      if (settings.blacklist.some(tag => tags.includes(tag.toLowerCase()))) {
        e.dataset.hidden = "true";
      }
    });
  }
};

const showIndexSettingsPopup = () => {
  if (document.getElementById("settings-popup")) {
    return;
  }

  const popup = document.createElement("div");
  popup.id = "settings-popup";

  const overlay = document.createElement("div");
  overlay.addEventListener("click", () => popup.remove());
  overlay.classList.add("settings-overlay");

  const content = document.createElement("div");
  content.classList.add("settings-content");

  const header = document.createElement("div");
  header.classList.add("settings-header");

  const title = document.createElement("span");
  title.textContent = "Settings";
  header.appendChild(title);

  const container = document.createElement("div");
  container.classList.add("settings-container");
  {
    const section = document.createElement("div");
    section.classList.add("settings-section", "blacklist");

    const label = document.createElement("span");
    label.textContent = "Hide tags";

    const actions = document.createElement("div");
    const wrapper = document.createElement("div");
    wrapper.classList.add("wrapper");

    const input = document.createElement("input");
    input.type = "text";
    input.defaultValue = settings.blacklist.join(", ");
    input.addEventListener("input", () => {
      const blacklist = input.value
        .split(",")
        .map(e => e.trim())
        .filter(e => e);
      setSettings("blacklist", blacklist);
      applySettings();
    });

    wrapper.appendChild(input);
    actions.appendChild(wrapper);
    section.append(label, actions);
    container.appendChild(section);
  }

  const note = document.createElement("p");
  note.textContent = "Tags are case-insensitive, delimited by commas.";
  note.style.marginTop = "1rem";
  container.appendChild(note);

  const footer = document.createElement("div");
  footer.classList.add("footer");

  const closeBtn = document.createElement("button");
  closeBtn.classList.add("close");
  closeBtn.type = "button";
  closeBtn.textContent = "Close";
  closeBtn.addEventListener("click", () => popup.remove());

  footer.appendChild(closeBtn);
  container.appendChild(footer);
  content.append(header, container);
  popup.append(overlay, content);

  document.body.appendChild(popup);
};

const initIndexSettings = () => {
  const toggleBtn = document.querySelector(".toggle-filter");
  if (!toggleBtn) return;

  getSettings();
  applySettings();

  toggleBtn.addEventListener("click", () => showIndexSettingsPopup());
};

const init = () => {
  if ("serviceWorker" in navigator) {
    navigator.serviceWorker.getRegistrations().then(registrations => {
      for (let i = 0; i < registrations.length; i++) {
        registrations[i].unregister();
      }
    });
  }

  initIndexSettings();

  const reader = document.getElementById("reader");
  if (!reader) return;

  id = document.body.dataset.id;
  slug = document.body.dataset.slug;
  totalPages = Number(document.body.dataset.totalPages);

  currPageSpans = document.querySelectorAll(".current");
  currPageNum = Number(currPageSpans[0].textContent);
  pageStates.push(...Array.from({ length: totalPages }, (_, i) => ({ isViewing: i + 1 === currPageNum })));

  pageContainer = reader.querySelector(".page a");
  currPageImage = pageContainer.querySelector("img");
  ({ origin } = new URL(currPageImage.src));

  firstPageAnchors = document.querySelectorAll(".first");
  firstPageAnchors.forEach(e => {
    e.addEventListener("click", ev => {
      preventDefault(ev);
      if (currPageNum > 1) {
        changePage(1);
      }
    });
  });

  lastPageAnchors = document.querySelectorAll(".last");
  lastPageAnchors.forEach(e => {
    e.addEventListener("click", ev => {
      preventDefault(ev);
      if (currPageNum < totalPages) {
        changePage(totalPages);
      }
    });
  });

  prevPageAnchors = document.querySelectorAll(".prev");
  prevPageAnchors.forEach(e => {
    e.addEventListener("click", ev => {
      preventDefault(ev);
      if (currPageNum > 1) {
        changePage(currPageNum - 1);
      }
    });
  });

  nextPageAnchors = document.querySelectorAll(".next");
  nextPageAnchors.forEach(e => {
    e.addEventListener("click", ev => {
      preventDefault(ev);
      if (currPageNum < totalPages) {
        changePage(currPageNum + 1);
      }
    });
  });

  pageContainer.addEventListener("click", (ev: MouseEvent) => {
    preventDefault(ev);

    const target = ev.target as HTMLImageElement;
    const isPrev = ev.offsetX <= (target.clientWidth * readerSettings.zoomLevel) / 2;

    if (isPrev && currPageNum > 1) {
      changePage(currPageNum - 1);
    } else if (!isPrev && currPageNum < totalPages) {
      changePage(currPageNum + 1);
    }
  });

  let interval = 0;
  const scrollIntoView = () => {
    clearInterval(interval);
    interval = window.setInterval(() => {
      window.requestAnimationFrame(() => {
        if (!currPageImage.naturalHeight) return;
        clearInterval(interval);
        pageContainer.scrollIntoView({ block: "start", inline: "start" });
      });
    }, 10);
  };

  const observer = new MutationObserver(() => {
    scrollIntoView();
    initPreload();
  });
  observer.observe(pageContainer, { childList: true });

  window.addEventListener("keydown", ev => {
    if (currPageNum > 1 && (ev.code === "ArrowLeft" || ev.code === "KeyA" || ev.code === "KeyH")) {
      changePage(currPageNum - 1);
    } else if (currPageNum < totalPages && (ev.code === "ArrowRight" || ev.code === "KeyD" || ev.code === "KeyL")) {
      changePage(currPageNum + 1);
    }
  });

  getReaderSettings();
  applyReaderSettings();
  attachHandlers();

  scrollIntoView();
  initPreload();
};

if (document.readyState === "complete") {
  init();
} else {
  document.addEventListener("DOMContentLoaded", init);
}
