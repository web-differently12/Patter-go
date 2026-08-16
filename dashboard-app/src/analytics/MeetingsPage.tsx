import React from 'react';

export const MeetingsPage: React.FC = () => {
  return (
    <div className="space-y-6 text-slate-100 p-6 bg-[#08090a] min-h-screen font-sans">
      {/* Header */}
      <div className="border-b border-slate-800 pb-4">
        <h1 className="text-2xl font-bold tracking-tight text-white">Suivi des Réunions Visio & Copilote</h1>
        <p className="text-sm text-slate-400">Bots Zoom, Google Meet et Microsoft Teams : temps de parole, objections et suivi des actions post-réunion.</p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Talk Time & Diarization */}
        <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-5 shadow-lg space-y-4">
          <h3 className="font-bold text-white text-base">Temps de Parole & Diarisation</h3>
          <p className="text-xs text-slate-400">Répartition du temps d'expression durant les réunions visio.</p>

          <div className="space-y-3 text-xs">
            <div>
              <div className="flex justify-between font-medium mb-1">
                <span className="text-slate-300">Interlocuteurs (Prospects)</span>
                <span className="text-cyan-400 font-bold">62%</span>
              </div>
              <div className="w-full bg-slate-800 h-2 rounded-full overflow-hidden">
                <div className="bg-cyan-400 h-full" style={{ width: '62%' }}></div>
              </div>
            </div>

            <div>
              <div className="flex justify-between font-medium mb-1">
                <span className="text-slate-300">Agent IA EchoFlow</span>
                <span className="text-indigo-400 font-bold">38%</span>
              </div>
              <div className="w-full bg-slate-800 h-2 rounded-full overflow-hidden">
                <div className="bg-indigo-400 h-full" style={{ width: '38%' }}></div>
              </div>
            </div>
          </div>
        </div>

        {/* Objections & Copilot Suggestions */}
        <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-5 shadow-lg space-y-4">
          <h3 className="font-bold text-white text-base">Objections & Suggestions Copilote</h3>
          <p className="text-xs text-slate-400">Liste des points complexes soulevés et réponses guidées par l'IA.</p>

          <div className="space-y-3 text-xs">
            <div className="bg-slate-950 p-3 rounded-lg border border-slate-800 space-y-1">
              <div className="text-amber-300 font-semibold">"Pouvons-nous héberger les données On-Premise ?"</div>
              <div className="text-slate-400">Suggestion Copilote: "Proposer le déploiement hybride sur VPC dédié AWS/GCP."</div>
            </div>

            <div className="bg-slate-950 p-3 rounded-lg border border-slate-800 space-y-1">
              <div className="text-amber-300 font-semibold">"Quel est le délai de mise en production ?"</div>
              <div className="text-slate-400">Suggestion Copilote: "Indiquer un déploiement standard en 48h avec connecteur CRM."</div>
            </div>
          </div>
        </div>

        {/* Post-Meeting Action Items */}
        <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-5 shadow-lg space-y-4">
          <h3 className="font-bold text-white text-base">Suivi des Actions Post-Réunion</h3>
          <p className="text-xs text-slate-400">Taux d'exécution des compte-rendus et tâches attribuées.</p>

          <div className="space-y-3 text-xs">
            <div className="flex justify-between items-center bg-slate-950 p-3 rounded-lg border border-slate-800">
              <span className="text-slate-300">Compte-Rendus Envoyés Automatiquement</span>
              <span className="text-emerald-400 font-bold bg-emerald-500/10 px-2 py-0.5 rounded border border-emerald-500/20">98.5%</span>
            </div>

            <div className="flex justify-between items-center bg-slate-950 p-3 rounded-lg border border-slate-800">
              <span className="text-slate-300">Tâches CRM / Follow-up créées</span>
              <span className="text-cyan-400 font-bold bg-cyan-500/10 px-2 py-0.5 rounded border border-cyan-500/20">142 tâches</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
