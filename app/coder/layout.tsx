import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Coder",
  description: "Collaborative IDE Platform by Coder",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return children;
}
