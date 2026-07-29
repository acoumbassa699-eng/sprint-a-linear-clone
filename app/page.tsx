import { Hero3DStage } from "@/components/hero-3d-stage"
import { auth } from "@/lib/auth";
import { redirect } from "next/navigation";
import Link from "next/link";
import { signOut } from "@/lib/auth";

export default async function Home() {
  const session = await auth();

  if (!session) {
    redirect("/login");
  }

  return (
    <main>
      {/* Header with user info */}
      <header className="fixed top-0 left-0 right-0 z-50 bg-zinc-900/80 backdrop-blur-xl border-b border-zinc-800">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            <div className="flex items-center gap-4">
              <h1 className="text-xl font-bold text-white">Axon</h1>
              <span className="text-zinc-500 text-sm">Linear Clone</span>
            </div>
            <div className="flex items-center gap-4">
              <div className="flex items-center gap-3">
                <div className="w-8 h-8 bg-violet-600 rounded-full flex items-center justify-center text-white text-sm font-medium">
                  {session.user?.name?.[0] || session.user?.email?.[0] || "U"}
                </div>
                <div className="hidden sm:block">
                  <p className="text-sm text-white">{session.user?.name || "User"}</p>
                  <p className="text-xs text-zinc-500">{session.user?.email}</p>
                </div>
              </div>
              <button
                onClick={() => signOut({ callbackUrl: "/login" })}
                className="px-4 py-2 text-sm text-zinc-300 hover:text-white hover:bg-zinc-800 rounded-lg transition-colors"
              >
                Déconnexion
              </button>
            </div>
          </div>
        </div>
      </header>

      {/* Hero Section */}
      <div className="pt-24 pb-16 px-4">
        <div className="max-w-4xl mx-auto text-center">
          <h2 className="text-4xl sm:text-5xl font-bold text-white mb-6">
            Bienvenue sur <span className="text-violet-400">Axon</span>
          </h2>
          <p className="text-xl text-zinc-400 mb-8">
            Plateforme de développement pour équipes natives AI
          </p>
          <div className="flex flex-wrap justify-center gap-4">
            <Link
              href="#features"
              className="px-6 py-3 bg-violet-600 hover:bg-violet-500 text-white font-medium rounded-lg transition-colors"
            >
              Découvrir les fonctionnalités
            </Link>
            <button className="px-6 py-3 bg-zinc-800 hover:bg-zinc-700 text-white font-medium rounded-lg transition-colors border border-zinc-700">
              Voir la documentation
            </button>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <Hero3DStage />
    </main>
  );
}
