export enum Mode {
  Normal,
  Strip
}

export interface Settings {
  blacklist: string[];
}

export const settings: Settings = {
  blacklist: ["yaoi", "crossdressing"]
};

export const saveSettings = () => {
  localStorage.setItem("indexSettings", JSON.stringify(settings));
};

export const getSettings = () => {
  const localSettings = localStorage.getItem("indexSettings");
  if (localSettings) {
    const obj = JSON.parse(localSettings);
    Object.assign(settings, obj);
  } else saveSettings();
};

export const setSettings = (key: string, value: any) => {
  if (!(key in settings)) {
    return;
  }

  settings[key] = value;
  saveSettings();
};

export interface ReaderSettings {
  mode?: Mode;
  maxWidth?: number;
  zoomLevel?: number;
}

export const readerSettings: ReaderSettings = {
  mode: 0,
  maxWidth: 1366,
  zoomLevel: 1.0
};

export const saveReaderSettings = () => {
  localStorage.setItem("settings", JSON.stringify(readerSettings));
};

export const getReaderSettings = () => {
  const localSettings = localStorage.getItem("settings");
  if (localSettings) {
    const obj = JSON.parse(localSettings);
    Object.assign(readerSettings, obj);
  } else saveReaderSettings();
};

export const setReaderSettings = (key: string, value: any) => {
  if (!(key in readerSettings) || (key === "zoomLevel" && (value < 0.1 || value > 5.0))) {
    return;
  }

  readerSettings[key] = value;
  saveReaderSettings();
};
