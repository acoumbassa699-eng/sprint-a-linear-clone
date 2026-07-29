import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Optimus IDE Collab",
  description: "Collaborative IDE Platform",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return children;
}
