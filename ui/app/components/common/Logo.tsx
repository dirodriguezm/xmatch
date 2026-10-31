"use client";

import Image from "next/image";

interface LogoProps {
  size?: "small" | "default" | "large" | "xlarge";
}

const sizeMap = {
  small: 24,
  default: 32,
  large: 48,
  xlarge: 80,
};

export function Logo({ size = "default" }: LogoProps) {
  const dimension = sizeMap[size];

  return (
    <Image
      src="/xwave-icon.svg"
      alt="XWave Logo"
      width={dimension}
      height={dimension}
      priority
    />
  );
}
