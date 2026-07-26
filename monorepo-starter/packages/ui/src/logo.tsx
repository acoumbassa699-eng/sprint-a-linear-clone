import type { SVGProps } from "react"

export function AxonMark(props: SVGProps<SVGSVGElement>) {
  return (
    <svg
      viewBox="0 0 32 32"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      role="img"
      aria-label="Axon logo"
      {...props}
    >
      {/* Rounded badge */}
      <rect x="1" y="1" width="30" height="30" rx="8" className="fill-transparent" />
      <rect
        x="1"
        y="1"
        width="30"
        height="30"
        rx="8"
        stroke="currentColor"
        strokeOpacity="0.25"
        strokeWidth="1.5"
      />
      {/* Two signal paths converging into an "A" */}
      <path
        d="M9 23 L16 9 L23 23"
        stroke="currentColor"
        strokeWidth="2.4"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      {/* Synapse connector */}
      <path d="M12 18 H20" stroke="currentColor" strokeWidth="2.4" strokeLinecap="round" />
      {/* Active signal node */}
      <circle cx="16" cy="9" r="2.6" className="fill-[color:var(--axon-accent,#6366f1)]" />
    </svg>
  )
}

export function Logo({ className }: { className?: string }) {
  return (
    <div className={`flex items-center gap-2 ${className ?? ""}`}>
      <AxonMark className="h-6 w-6 text-white" />
      <span className="text-white font-semibold tracking-tight">Axon</span>
    </div>
  )
}
