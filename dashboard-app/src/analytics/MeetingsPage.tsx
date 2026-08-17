import React, { useState } from 'react';

export const MeetingsPage: React.FC = () => {
  // Connection & Video Render
  const [videoResolution, setVideoResolution] = useState<'720p' | '1080p'>('1080p');
  const [botMode, setBotMode] = useState<'rtms_native' | 'signed_in_bot'>('signed_in_bot');
  const [multiTrackIsolation, setMultiTrackIsolation] = useState<boolean>(true);

  // Advanced Bot Customization & Behavior
  const [botName, setBotName] = useState('Patter Meeting Assistant');
  const [botAvatarUrl, setBotAvatarUrl] = useState('https://kallflow.ai/assets/avatar-bot.png');
  const [joinTiming, setJoinTiming] = useState<'at_start' | '2_min_before' | 'on_call_start'>('at_start');
  const [recordingConsentPrompt, setRecordingConsentPrompt] = useState<boolean>(true);
  const [captureMode, setCaptureMode] = useState<'audio_video' | 'audio_only'>('audio_video');

  // Integrations & Realtime Automation
  const [calendarSync, setCalendarSync] = useState<boolean>(true);
  const [chatBroadcasting, setChatBroadcasting] = useState<boolean>(true);
  const [liveWebhooks, setLiveWebhooks] = useState<boolean>(true);
  const [autoActionItems, setAutoActionItems] = useState<boolean>(true);

  return (
    <div className="space-y-6 text-slate-100 p-6 bg-[#08090a] min-h-screen font-sans">
      {/* Header */}
      <div className="border-b border-slate-800 pb-4 flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-white">Suivi des Reunions Visio & Moteur Bot Native</h1>
          <p className="text-sm text-slate-400">Bots Zoom (RTMS), Google Meet & Microsoft Teams: capture audio/video isolee, diarisation, chat et copilote en direct.</p>
        </div>
        <span className="text-xs font-mono bg-cyan-950 text-cyan-400 border border-cyan-800 px-3 py-1 rounded-full font-semibold self-start md:self-auto">
          Moteur Visio Native v2.4
        </span>
      </div>

      {/* Main Bot Capabilities & Configuration Panels */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Core Video & Media Capture Configuration */}
        <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-5 shadow-lg space-y-4">
          <h3 className="font-bold text-white text-base">Configuration des Flux Video & Audio Multi-Pistes</h3>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 text-xs">
            <div className="bg-slate-950 p-3 rounded-lg border border-slate-800 space-y-1.5">
              <label className="font-semibold text-slate-300 block">Mode de Connexion Bot</label>
              <select
                value={botMode}
                onChange={(e) => setBotMode(e.target.value as any)}
                className="w-full bg-slate-900 border border-slate-700 rounded p-2 text-slate-200 focus:outline-none focus:border-cyan-500"
              >
                <option value="signed_in_bot">Bot Identifie (Zoom / Teams / Meet)</option>
                <option value="rtms_native">Flux Direct RTMS / Native Media API</option>
              </select>
              <p className="text-[11px] text-slate-500">Capture le flux audio et video directement sous votre propre marque.</p>
            </div>

            <div className="bg-slate-950 p-3 rounded-lg border border-slate-800 space-y-1.5">
              <label className="font-semibold text-slate-300 block">Resolution Video Render</label>
              <select
                value={videoResolution}
                onChange={(e) => setVideoResolution(e.target.value as any)}
                className="w-full bg-slate-900 border border-slate-700 rounded p-2 text-slate-200 focus:outline-none focus:border-cyan-500"
              >
                <option value="1080p">1080p Full HD (Haute Fidelite)</option>
                <option value="720p">720p HD (Economique)</option>
              </select>
              <p className="text-[11px] text-slate-500">Qualite de rendu video pour la retransmission et l'analyse visuelle.</p>
            </div>

            <div className="bg-slate-950 p-3 rounded-lg border border-slate-800 space-y-1.5">
              <label className="font-semibold text-slate-300 block">Mode de Capture Medias</label>
              <select
                value={captureMode}
                onChange={(e) => setCaptureMode(e.target.value as any)}
                className="w-full bg-slate-900 border border-slate-700 rounded p-2 text-slate-200 focus:outline-none focus:border-cyan-500"
              >
                <option value="audio_video">Audio + Video Multi-Pistes</option>
                <option value="audio_only">Audio Seul (Transcription Discrete)</option>
              </select>
            </div>

            <div className="bg-slate-950 p-3 rounded-lg border border-slate-800 space-y-1.5">
              <label className="font-semibold text-slate-300 block">Isolation Audio Multi-Pistes</label>
              <div className="flex items-center justify-between pt-1">
                <span className="text-slate-400">Pistes PCM Separees</span>
                <input
                  type="checkbox"
                  checked={multiTrackIsolation}
                  onChange={(e) => setMultiTrackIsolation(e.target.checked)}
                  className="h-4 w-4 rounded bg-slate-900 border-slate-700 text-cyan-500 focus:ring-cyan-500"
                />
              </div>
              <p className="text-[11px] text-slate-500">Isole le flux micro de chaque participant pour une diarisation parfaite.</p>
            </div>
          </div>
        </div>

        {/* Identity, Calendar Sync & Realtime Automation */}
        <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-5 shadow-lg space-y-4">
          <h3 className="font-bold text-white text-base">Identite du Bot, Calendriers & Automatisations</h3>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 text-xs font-mono">
            <div>
              <label className="text-slate-300 block mb-1">Nom Affiche du Bot Visio</label>
              <input
                type="text"
                value={botName}
                onChange={(e) => setBotName(e.target.value)}
                className="w-full bg-slate-950 border border-slate-800 rounded-lg p-2 text-slate-200"
              />
            </div>

            <div>
              <label className="text-slate-300 block mb-1">Timing de Connexion Bot</label>
              <select
                value={joinTiming}
                onChange={(e) => setJoinTiming(e.target.value as any)}
                className="w-full bg-slate-950 border border-slate-800 rounded-lg p-2 text-slate-200"
              >
                <option value="at_start">Au Debut de la Reunion</option>
                <option value="2_min_before">2 Minutes Avant</option>
                <option value="on_call_start">A la Premiere Entree Participant</option>
              </select>
            </div>

            <div className="sm:col-span-2 grid grid-cols-2 gap-2 text-[11px] pt-1">
              <label className="flex items-center gap-2 p-2 bg-slate-950 rounded border border-slate-800">
                <input
                  type="checkbox"
                  checked={recordingConsentPrompt}
                  onChange={(e) => setRecordingConsentPrompt(e.target.checked)}
                  className="accent-cyan-500"
                />
                <span>Annonce de Consentement Enregistrement</span>
              </label>

              <label className="flex items-center gap-2 p-2 bg-slate-950 rounded border border-slate-800">
                <input
                  type="checkbox"
                  checked={calendarSync}
                  onChange={(e) => setCalendarSync(e.target.checked)}
                  className="accent-cyan-500"
                />
                <span>Auto-Sync Calendriers (Google/Outlook V2)</span>
              </label>

              <label className="flex items-center gap-2 p-2 bg-slate-950 rounded border border-slate-800">
                <input
                  type="checkbox"
                  checked={chatBroadcasting}
                  onChange={(e) => setChatBroadcasting(e.target.checked)}
                  className="accent-cyan-500"
                />
                <span>Diffusion de Réponses dans le Chat Visio</span>
              </label>

              <label className="flex items-center gap-2 p-2 bg-slate-950 rounded border border-slate-800">
                <input
                  type="checkbox"
                  checked={autoActionItems}
                  onChange={(e) => setAutoActionItems(e.target.checked)}
                  className="accent-cyan-500"
                />
                <span>Extraction Automatique des Taches CRM</span>
              </label>
            </div>
          </div>
        </div>
      </div>

      {/* Analytics Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Talk Time & Diarization */}
        <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-5 shadow-lg space-y-4">
          <h3 className="font-bold text-white text-base">Temps de Parole & Diarisation Horodatée</h3>
          <p className="text-xs text-slate-400">Repartition exacte du temps d'expression par locuteur.</p>

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
                <span className="text-slate-300">Agent IA Patter Engine</span>
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
          <h3 className="font-bold text-white text-base">Objections & Suggestions Copilote Temps Réel</h3>
          <p className="text-xs text-slate-400">Analyse sémantique des questions complexes et réponses recommandées.</p>

          <div className="space-y-3 text-xs">
            <div className="bg-slate-950 p-3 rounded-lg border border-slate-800 space-y-1">
              <div className="text-amber-300 font-semibold">"Pouvons-nous herberger les donnees On-Premise ?"</div>
              <div className="text-slate-400">Suggestion Copilote: "Proposer le deploiement hybride sur VPC dedie AWS/GCP."</div>
            </div>

            <div className="bg-slate-950 p-3 rounded-lg border border-slate-800 space-y-1">
              <div className="text-amber-300 font-semibold">"Quel est le delai de mise en production ?"</div>
              <div className="text-slate-400">Suggestion Copilote: "Indiquer un deploiement standard en 48h avec connecteur CRM."</div>
            </div>
          </div>
        </div>

        {/* Post-Meeting Action Items */}
        <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-5 shadow-lg space-y-4">
          <h3 className="font-bold text-white text-base">Suivi des Actions Post-Réunion & Synchro CRM</h3>
          <p className="text-xs text-slate-400">Exportation vers Hubspot, Salesforce et Webhooks.</p>

          <div className="space-y-3 text-xs">
            <div className="flex justify-between items-center bg-slate-950 p-3 rounded-lg border border-slate-800">
              <span className="text-slate-300">Compte-Rendus Envoyes Automatiquement</span>
              <span className="text-emerald-400 font-bold bg-emerald-500/10 px-2 py-0.5 rounded border border-emerald-500/20">98.5%</span>
            </div>

            <div className="flex justify-between items-center bg-slate-950 p-3 rounded-lg border border-slate-800">
              <span className="text-slate-300">Taches CRM / Follow-up creees</span>
              <span className="text-cyan-400 font-bold bg-cyan-500/10 px-2 py-0.5 rounded border border-cyan-500/20">142 taches</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
