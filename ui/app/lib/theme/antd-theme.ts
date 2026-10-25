import { theme, ThemeConfig } from "antd";

export const darkTheme: ThemeConfig = {
  algorithm: theme.darkAlgorithm,
  token: {
    // Primary colors (space/astronomy theme)
    colorPrimary: "#3768e8",
    colorInfo: "#3768e8",
    colorLink: "#7ea2ff",
    colorLinkHover: "#a9c1ff",
    colorBgContainer: "#12151c",
    colorBgElevated: "#1a1e27",
    colorBgLayout: "#0b0d12",
    colorBorder: "#272c37",
    colorBorderSecondary: "#1a1e27",

    // Text colors
    colorText: "#e4e7ec",
    colorTextSecondary: "#a3a9b5",
    colorTextTertiary: "#6f7684",

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
      headerBg: "#12151c",
      siderBg: "#12151c",
      bodyBg: "#0b0d12",
      triggerBg: "#1a1e27",
    },
    Input: {
      activeBg: "#1a1e27",
      addonBg: "#1a1e27",
    },
    Select: {
      optionSelectedBg: "#1a1e27",
    },
    Table: {
      headerBg: "#1a1e27",
      rowHoverBg: "#1a1e27",
    },
    Button: {
      primaryShadow: "none",
    },
    Card: {
      colorBgContainer: "#12151c",
    },
  },
};
