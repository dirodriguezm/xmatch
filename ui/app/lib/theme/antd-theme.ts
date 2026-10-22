import { theme, ThemeConfig } from "antd";

export const darkTheme: ThemeConfig = {
  algorithm: theme.darkAlgorithm,
  token: {
    // Primary colors (space/astronomy theme)
    colorPrimary: "#1677ff",
    colorBgContainer: "#141414",
    colorBgElevated: "#1f1f1f",
    colorBgLayout: "#0a0a0a",
    colorBorder: "#303030",
    colorBorderSecondary: "#1f1f1f",

    // Text colors
    colorText: "#e6e6e6",
    colorTextSecondary: "#a6a6a6",
    colorTextTertiary: "#666666",

    // Component specific
    borderRadius: 6,
    fontFamily:
      "var(--font-geist-sans), -apple-system, BlinkMacSystemFont, sans-serif",

    // Status colors
    colorSuccess: "#52c41a",
    colorWarning: "#faad14",
    colorError: "#ff4d4f",
  },
  components: {
    Layout: {
      headerBg: "#141414",
      siderBg: "#141414",
      bodyBg: "#0a0a0a",
      triggerBg: "#1f1f1f",
    },
    Input: {
      activeBg: "#1f1f1f",
      addonBg: "#1f1f1f",
    },
    Select: {
      optionSelectedBg: "#1f1f1f",
    },
    Table: {
      headerBg: "#1f1f1f",
      rowHoverBg: "#1f1f1f",
    },
    Button: {
      primaryShadow: "none",
    },
    Card: {
      colorBgContainer: "#141414",
    },
  },
};
