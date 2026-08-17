import React from 'react';

export const AdminSettingsPage: React.FC = () => {
  return (
    <div className="space-y-6 text-slate-100 p-6 bg-[#08090a] min-h-screen font-sans">
      <div className="border-b border-zinc-800 pb-4">
        <span className="text-xs font-mono font-bold text-cyan-400">BetterAuth Security & System Config</span>
        <h1 className="text-2xl font-bold tracking-tight text-white mt-0.5">Paramètres d'Administration & Sécurité</h1>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 text-xs">
        <div className="bg-zinc-950 border border-zinc-800 rounded-2xl p-5 space-y-3">
          <h3 className="font-bold text-white text-sm">Authentification & 2FA (BetterAuth)</h3>
          <p className="text-zinc-400">Authentification Google OAuth, Magic Links et 2FA/TOTP activees.</p>
          <div className="text-emerald-400 font-mono">Sessions JWT Bridge Go: ACTIF</div>
        </div>

        <div className="bg-zinc-950 border border-zinc-800 rounded-2xl p-5 space-y-3">
          <h3 className="font-bold text-white text-sm">Fournisseur d'E-mails Transactionnels</h3>
          <p className="text-zinc-400">Configuration SMTP Mailgun / API Resend pour notifications.</p>
          <div className="text-zinc-300 font-mono">API Resend: Connecte</div>
        </div>
      </div>
    </div>
  );
};
