

/* ═══════════ خان — Khan Chat Frontend (Vue 3) v1.0.3 ═══════════
   Full-featured: public/private rooms, offline msgs, presence, typing,
   read receipts, edit/delete, mentions, search, forward/reply/pin,
   polls, urgent, unread counters, desktop notif, Persian calendar,
   drag&drop files, avatar upload, stickers, departments, backup, logs
*/
const { createApp } = Vue;

// ─── دوزبانه: فارسی / English ───
const I18N = {
  fa: {
    // Login / Setup
    loginSubtitle: 'چت سازمانی سبک برای شبکه محلی شما',
    username: 'نام کاربری',
    usernamePlaceholder: 'نام کاربری خود را وارد کنید',
    password: 'رمز عبور',
    passwordPlaceholder: '••••••••',
    loggingIn: 'در حال ورود...',
    login: 'ورود',
    setupAdmin: 'نام کاربری مدیر',
    setupDisplayName: 'نام نمایشی',
    setupPassword: 'رمز عبور',
    setupCreate: '🚀 ساخت مدیر',
    // Sidebar
    searchChats: 'جستجوی گفتگوها...',
    chats: 'گفتگوها',
    rooms: 'اتاق‌ها',
    search: 'جستجو',
    noChats: 'هنوز گفتگویی نیست',
    noChatsDesc: 'از تب اتاق‌ها یک اتاق بسازید یا به اتاق عمومی بپیوندید',
    publicRooms: 'اتاق‌های عمومی',
    privateRooms: 'اتاق‌های خصوصی',
    clickToJoin: 'برای عضویت کلیک کنید',
    createRoom: 'ساخت اتاق جدید',
    newRoom: 'اتاق جدید',
    createRoomDesc: 'گروه، عمومی، خصوصی یا کانال بسازید',
    smartSearch: 'جستجوی هوشمند...',
    searchMessages: 'جستجو در پیام‌ها',
    adminPanel: 'پنل مدیریت',
    logout: 'خروج',
    online: 'آنلاین',
    noOnline: 'هیچ‌کس آنلاین نیست',
    allUsers: 'همه کاربران',
    admin: 'مدیر',
    offline: 'آفلاین',
    members: 'عضو',
    // Chat
    khanChat: 'خان',
    urgentMessage: 'پیام فوری در این اتاق وجود دارد',
    loadOlder: 'پیام‌های قدیمی‌تر',
    forwarded: 'فوروارد شده',
    urgent: 'فوری',
    edited: 'ویرایش شده',
    votes: 'رای',
    pollClosed: 'بسته شده',
    viewProfile: 'مشاهده پروفایل',
    startChat: 'شروع چت',
    searchUsers: 'جستجوی کاربران',
    noUsersFound: 'کاربری یافت نشد',
    searchResults: 'نتایج جستجو',
    close: 'بستن',
    status: 'وضعیت',
    welcome: 'به خان خوش آمدید',
    welcomeDesc: 'یک اتاق را از سمت راست انتخاب کنید',
    noMessages: 'هنوز پیامی نیست',
    noMessagesDesc: 'اولین پیام را بفرستید!',
    typing: 'در حال نوشتن',
    someone: 'یک نفر',
    replyingTo: 'پاسخ به',
    emoji: 'ایموجی',
    stickers: 'استیکر',
    messagePlaceholder: 'پیام خود را بنویسید...',
    attach: 'پیوست فایل',
    dropFile: 'فایل را رها کنید تا ارسال شود',
    // Room info
    roomInfo: 'اطلاعات اتاق',
    removeMember: 'حذف عضو',
    addMember: 'افزودن عضو',
    selectUser: 'انتخاب کاربر...',
    pinnedMessages: 'پیام‌های پین شده',
    noPinned: 'پیامی پین نشده',
    // Admin
    close: 'بستن',
    users: 'کاربران',
    license: 'لایسنس',
    network: 'شبکه',
    backup: 'پشتیبان‌گیری',
    logs: 'لاگ‌ها',
    departments: 'بخش‌ها',
    userManagement: 'مدیریت کاربران',
    newUser: 'کاربر جدید',
    user: 'کاربر',
    role: 'نقش',
    status: 'وضعیت',
    actions: 'عملیات',
    demote: 'کاهش نقش',
    promote: 'افزایش نقش',
    maxUsers: 'حداکثر کاربران',
    company: 'سازمان',
    expires: 'انقضا',
    licenseWarning: 'لایسنس دستکاری شده! محدودیت به ۵ کاربر کاهش یافت.',
    applyLicense: 'اعمال لایسنس',
    removeLicense: 'حذف لایسنس',
    serverAddress: 'آدرس سرور',
    port: 'پورت',
    version: 'نسخه',
    addressType: 'نوع آدرس',
    address: 'آدرس',
    ipPlaceholder: 'آدرس IP یا DNS سرور',
    saveSettings: 'ذخیره تنظیمات',
    backupNow: 'پشتیبان‌گیری الان',
    backupFile: 'فایل پشتیبان',
    noBackups: 'هنوز پشتیبانی ساخته نشده',
    restore: 'بازیابی',
    restoreConfirm: 'بازیابی این نسخه پشتیبان؟ همه داده فعلی جایگزین می‌شود.',
    restored: 'بازیابی انجام شد',
    exportExcel: 'خروجی اکسل',
    exportDone: 'خروجی کاربران آماده شد',
    newDepartment: 'بخش جدید',
    myProfile: 'پروفایل من',
    displayName: 'نام نمایشی',
    avatarUpload: 'آپلود آواتار',
    save: 'ذخیره',
    changePassword: 'تغییر رمز عبور',
    currentPassword: 'رمز فعلی',
    newPassword: 'رمز جدید',
    // Poll
    createPoll: 'نظرسنجی جدید',
    pollQuestion: 'سوال نظرسنجی',
    pollOptions: 'گزینه‌ها',
    option: 'گزینه',
    addOption: 'افزودن گزینه',
    create: 'ساخت',
    cancel: 'انصراف',
    // Message menu
    forwardTo: 'فوروارد به...',
    reply: 'پاسخ',
    forward: 'فوروارد',
    pin: 'پین',
    unpin: 'برداشتن پین',
    delete: 'حذف',
    noResults: 'نتیجه‌ای یافت نشد',
    searchInRoom: 'جستجو در پیام‌ها...',
    // Rooms
    roomName: 'نام اتاق',
    roomType: 'نوع اتاق',
    groupRoom: 'گروه',
    publicRoom: 'عمومی',
    privateRoom: 'خصوصی (دعوتی)',
    channelRoom: 'کانال',
    department: 'بخش',
    none: 'بدون بخش',
    // Misc
    online: 'آنلاین',
    noOnline: 'هیچ‌کس آنلاین نیست',
    allUsers: 'همه کاربران',
    admin: 'مدیر',
    offline: 'آفلاین',
    lastSeen: 'آخرین بازدید',
  },
  en: {
    loginSubtitle: 'Lightweight team chat for your local network',
    username: 'Username',
    usernamePlaceholder: 'Enter your username',
    password: 'Password',
    passwordPlaceholder: '••••••••',
    loggingIn: 'Logging in...',
    login: 'Login',
    setupAdmin: 'Admin username',
    setupDisplayName: 'Display name',
    setupPassword: 'Password',
    setupCreate: '🚀 Create Admin',
    searchChats: 'Search chats...',
    chats: 'Chats',
    rooms: 'Rooms',
    search: 'Search',
    noChats: 'No chats yet',
    noChatsDesc: 'Create a room from the Rooms tab or join a public room',
    publicRooms: 'Public rooms',
    privateRooms: 'Private rooms',
    clickToJoin: 'Click to join',
    createRoom: 'Create new room',
    newRoom: 'New room',
    createRoomDesc: 'Create group, public, private or channel',
    smartSearch: 'Smart search...',
    searchMessages: 'Search messages',
    adminPanel: 'Admin Panel',
    logout: 'Logout',
    online: 'Online',
    noOnline: 'No one is online',
    allUsers: 'All users',
    admin: 'Admin',
    offline: 'Offline',
    members: 'members',
    khanChat: 'Khan',
    urgentMessage: 'Urgent message in this room',
    loadOlder: 'Load older messages',
    forwarded: 'Forwarded',
    urgent: 'URGENT',
    edited: 'edited',
    votes: 'votes',
    pollClosed: 'closed',
    viewProfile: 'View profile',
    startChat: 'Start chat',
    searchUsers: 'Search users',
    noUsersFound: 'No users found',
    searchResults: 'Search results',
    close: 'Close',
    status: 'Status',
    welcome: 'Welcome to Khan',
    welcomeDesc: 'Select a room from the left',
    noMessages: 'No messages yet',
    noMessagesDesc: 'Send the first message!',
    typing: 'is typing',
    someone: 'someone',
    replyingTo: 'Reply to',
    emoji: 'Emoji',
    stickers: 'Stickers',
    messagePlaceholder: 'Type a message...',
    attach: 'Attach file',
    dropFile: 'Drop file to send',
    roomInfo: 'Room info',
    removeMember: 'Remove member',
    addMember: 'Add member',
    selectUser: 'Select user...',
    pinnedMessages: 'Pinned messages',
    noPinned: 'No pinned messages',
    close: 'Close',
    users: 'Users',
    license: 'License',
    network: 'Network',
    backup: 'Backup',
    logs: 'Logs',
    departments: 'Departments',
    userManagement: 'User Management',
    newUser: 'New User',
    user: 'User',
    role: 'Role',
    status: 'Status',
    actions: 'Actions',
    demote: 'Demote',
    promote: 'Promote',
    maxUsers: 'Max users',
    company: 'Company',
    expires: 'Expires',
    licenseWarning: 'License tampered! Limit reduced to 5 users.',
    applyLicense: 'Apply License',
    removeLicense: 'Remove License',
    serverAddress: 'Server address',
    port: 'Port',
    version: 'Version',
    addressType: 'Address type',
    address: 'Address',
    ipPlaceholder: 'Server IP or DNS address',
    saveSettings: 'Save Settings',
    backupNow: 'Backup Now',
    backupFile: 'Backup file',
    noBackups: 'No backups yet',
    restore: 'Restore',
    restoreConfirm: 'Restore this backup? All current data will be replaced.',
    restored: 'Restore complete',
    exportExcel: 'Export Excel',
    exportDone: 'Users exported',
    newDepartment: 'New Department',
    myProfile: 'My Profile',
    displayName: 'Display name',
    avatarUpload: 'Upload avatar',
    save: 'Save',
    changePassword: 'Change Password',
    currentPassword: 'Current password',
    newPassword: 'New password',
    createPoll: 'New Poll',
    pollQuestion: 'Poll question',
    pollOptions: 'Options',
    option: 'Option',
    addOption: 'Add option',
    create: 'Create',
    cancel: 'Cancel',
    forwardTo: 'Forward to...',
    reply: 'Reply',
    forward: 'Forward',
    pin: 'Pin',
    unpin: 'Unpin',
    delete: 'Delete',
    noResults: 'No results',
    searchInRoom: 'Search in room...',
    roomName: 'Room name',
    roomType: 'Room type',
    groupRoom: 'Group',
    publicRoom: 'Public',
    privateRoom: 'Private (invite)',
    channelRoom: 'Channel',
    department: 'Department',
    none: 'None',
    lastSeen: 'Last seen',
  },
};

const EMOJIS = ['😀','😁','😂','🤣','😊','😍','😘','😎','🤔','🙄','😅','😉','👍','👎','👏','🙏','🔥','❤️','💔','💯','🎉','🎊','😢','😭','🤝','✅','❌','⭐','🚀','💪','🤝','👌','🙌','🤞','✌️','💡','🎯','📌','⚠️','❓','❗'];
const STICKERS = ['😀','😂','😍','🔥','👍','👏','🎉','😢','🤔','😎','🥳','😴','🤯','😭','🙏','💪','❤️','💯','⭐','🚀','🍕','☕','🌙','☀️'];

const app = createApp({
  data() {
    return {
      // Language & Theme
      uiLang: localStorage.getItem('khan_lang') || 'fa',
      uiTheme: localStorage.getItem('khan_theme') || 'khan-dark',
      themeList: [
        { id:'khan',     label:'خان',    swatch:'linear-gradient(135deg,#e8a33d,#171410)' },
        { id:'sapphire', label:'یاقوت',  swatch:'linear-gradient(135deg,#3b82f6,#0f131c)' },
        { id:'spring',   label:'بهار',   swatch:'linear-gradient(135deg,#10b981,#0e1912)' },
        { id:'night',    label:'شب',     swatch:'linear-gradient(135deg,#8b5cf6,#151022)' },
        { id:'sunset',   label:'غروب',   swatch:'linear-gradient(135deg,#f97316,#1a110c)' },
        { id:'frost',    label:'برف',    swatch:'linear-gradient(135deg,#94a3b8,#10141b)' },
      ],
      i18n: I18N,
      // Server & session
      serverInfo: null,
      session: { token: localStorage.getItem('khan_token') || '', user: null },
      setupMode: false,
      setupForm: { admin_username: '', admin_password: '', company_name: '' },
      setupBusy: false,
      setupError: '',
      loginForm: { username: '', password: '' },
      loginBusy: false,
      loginError: '',

      // Data
      users: [],
      rooms: [],
      messages: {},
      roomMembers: {},
      currentRoom: null,
      currentRoomId: null,
      hasMore: {},
      loadingMore: false,
      departments: [],
      publicRooms: [],
      privateRooms: [],

      // UI state
      chatSearch: '',
      sidebarTab: 'chats',
      mobileSidebar: false,
      mobileChatOpen: false,
      showRoomInfo: false,
      showProfile: false,
      showUserProfile: false,
      selectedUser: null,
      userSearchQuery: '',
      userSearchResults: [],
      showChangePwd: false,
      pwdForm: { current: '', next: '' },
      showAdminPanel: false,
      adminTab: 'users',
      showCreateRoom: false,
      createRoomForm: { name: '', type: 'group', department: 0 },
      showCreateUser: false,
      createUserForm: { username: '', display_name: '', password: '', role: 'user' },
      showCreateDept: false,
      createDeptForm: { name: '' },
      license: { state: 'free', max_users: 20, licensed_to: '', expiry: '', error: '' },
      networkForm: { address_type: 'ip', ip: '', dns: '', port: 1727 },
      networkInfo: {},
      backups: [],
      serverLogs: '',
      inviteUserId: null,
      profileForm: { display_name: '' },

      // Unread
      readMap: {},
      unreadMap: {},
      notifEnabled: false,
      docTitle: document.title,

      // Composer
      draft: '',
      replyTarget: null,
      showStickers: false,
      showEmojis: false,
      editingMsg: null,
      uploadProgress: null,
      uploadingFile: false,
      dragOver: false,
      mentionPopupVisible: false,
      mentionCandidates: [],
      mentionSelected: null,
      mentionQuery: '',

      // Messages
      typingUsers: {},
      typingNames_: {},
      pinnedMessages: [],
      currentRoomUrgent: false,
      searchResults: [],
      showSearchOverlay: false,
      searchQuery: '',
      smartSearchQuery: '',
      contextMenu: null,
      forwardTarget: null,
      showForwardModal: false,
      showPollModal: false,
      pollForm: { question: '', options: ['', ''] },
      toasts: [],
      toastId: 0,
      deptCollapsed: {},

      // WS
      ws: null,
      wsRetry: 0,
      reconnectTimer: null,
      emojis: EMOJIS,
      stickers: STICKERS,
    };
  },

  computed: {
    filteredChats() {
      const q = this.chatSearch.toLowerCase();
      return this.rooms.filter(r => this.chatName(r).toLowerCase().includes(q));
    },

    totalUnread() {
      let t = 0;
      for (const k in this.unreadMap) t += this.unreadMap[k] || 0;
      return t;
    },

    currentRoomMembers() {
      if (!this.currentRoom) return [];
      return this.roomMembers[this.currentRoom.id] || [];
    },

    groupedMessages() {
      const msgs = this.messages[this.currentRoomId] || [];
      const groups = {};
      for (const m of msgs) {
        if (m.deleted_at) continue;
        const key = this.dayKey(m.created_at);
        if (!groups[key]) groups[key] = [];
        groups[key].push(m);
      }
      return groups;
    },

    typingNames() {
          const names = [];
          for (const id in this.typingUsers) {
            if (id == this.session.user.id) continue;
            const u = this.users.find(x => x.id == id);
            // Prefer the live WS-sent name, then the user record, then a generic label
            const stored = (this.typingNames_ || {})[id];
            names.push(stored || (u && (u.display_name || u.username)) || this.t('someone'));
          }
          return names;
        },

    canManageRoom() {
      const u = this.session.user;
      if (!u || !this.currentRoom) return false;
      if (['adm', 'sadm'].includes(u.role)) return true;
      if (u.role === 'sup') return this.currentRoom.type !== 'dm';
      return false;
    },

    isAdmin() {
      const r = this.session.user ? this.session.user.role : '';
      return ['adm', 'sadm'].includes(r);
    },

    // Alias so the template's `needsSetup` (v-if in login card) tracks setupMode
    needsSetup() {
      return this.setupMode;
    },

    // Users tab: only visible users (backend already excludes the hidden super admin)
    onlineVisibleUsers() {
      return (this.users || []).filter(u => u.id !== (this.session.user || {}).id && u.online);
    },

    visibleUsersSorted() {
      return (this.users || [])
        .filter(u => u.id !== (this.session.user || {}).id)
        .sort((a, b) => (b.online ? 1 : 0) - (a.online ? 1 : 0) || (a.display_name || a.username || '').localeCompare(b.display_name || b.username || ''));
    },

    licenseIcon() {
      return { free: '🆓', valid: '✅', tampered: '⚠️' }[this.license.state] || '🎫';
    },

    licenseTitle() {
      return { free: this.t('license') + ' — ' + this.t('maxUsers') + ' 20', valid: '✅ ' + this.t('applyLicense'), tampered: '⚠️' }[this.license.state] || '🎫';
    },

    dmPartnerOnline() {
      const p = this.dmPartner(this.currentRoom);
      return p ? this.isUserOnline(p) : false;
    },

    dmPartnerStatusText() {
      const p = this.dmPartner(this.currentRoom);
      if (!p) return '';
      if (this.isUserOnline(p)) return this.t('online');
      const ls = p.last_seen;
      if (ls) return this.t('lastSeen') + ': ' + this.relativeTime(ls);
      return this.t('offline');
    },

    roomHasOnline() {
      return this.currentRoomMembers.some(m => m.id !== this.session.user.id && this.isUserOnline(m));
    },

    roomOnlineCount() {
      return this.currentRoomMembers.filter(m => m.id !== this.session.user.id && this.isUserOnline(m)).length;
    },

    availableMembers() {
      const existing = this.currentRoomMembers.map(m => m.id);
      return this.users.filter(u => !existing.includes(u.id) && u.id !== this.session.user.id);
    },

    canPin() {
      return this.canManageRoom || (this.contextMenu && this.contextMenu.msg && this.contextMenu.msg.sender_id === this.session.user.id);
    },

    canUnpin() {
      return this.canManageRoom || (this.contextMenu && this.contextMenu.msg && this.contextMenu.msg.sender_id === this.session.user.id);
    },
  },

  methods: {
    /* ─────────── i18n ─────────── */
    t(key) {
      return (this.i18n[this.uiLang] && this.i18n[this.uiLang][key]) || this.i18n.fa[key] || key;
    },
    setLang(lang) {
      this.uiLang = lang;
      localStorage.setItem('khan_lang', lang);
      document.documentElement.lang = lang === 'fa' ? 'fa' : 'en';
      document.documentElement.dir = lang === 'fa' ? 'rtl' : 'ltr';
    },

    /* ─────────── Theme (v1.0.8: 6 bases × light/dark) ─────────── */
    get themeBase() {
      const t = this.uiTheme || 'dark';
      return t.replace(/-light$|-dark$/, '');
    },
    get themeMode() {
      const t = this.uiTheme || 'dark';
      return /-light$/.test(t) ? 'light' : 'dark';
    },
    applyTheme() {
      localStorage.setItem('khan_theme', this.uiTheme);
      document.documentElement.setAttribute('data-theme', this.uiTheme);
      const meta = document.querySelector('meta[name="theme-color"]');
      if (meta) {
        const cs = getComputedStyle(document.documentElement);
        const bg = cs.getPropertyValue('--bg').trim() || cs.getPropertyValue('--bg-deep').trim() || '#0e0f12';
        meta.content = bg;
      }
    },
    setThemeBase(base) {
      this.uiTheme = base + '-' + (this.themeMode === 'light' ? 'light' : 'dark');
      this.applyTheme();
    },
    toggleThemeMode() {
      this.uiTheme = this.themeBase + '-' + (this.themeMode === 'dark' ? 'light' : 'dark');
      this.applyTheme();
    },

    /* ─────────── API ─────────── */
    async api(path, options = {}) {
      const headers = { 'Content-Type': 'application/json' };
      if (this.session.token) headers['Authorization'] = 'Bearer ' + this.session.token;
      const res = await fetch(path, {
        ...options,
        headers: { ...headers, ...(options.headers || {}) },
      });
      const data = await res.json().catch(() => ({}));
      if (!res.ok) throw new Error(data.error || 'خطای سرور');
      return data;
    },

    /* ─────────── Setup / Login ─────────── */
    async checkSetup() {
      try {
        const data = await this.api('/api/setup/needs-setup');
        this.setupMode = data.needs_setup;
        if (!this.setupMode && this.session.token) {
          await this.loadSession();
        }
      } catch (e) {
        this.setupMode = true;
      }
    },

    async loadSession() {
      try {
        const me = await this.api('/api/auth/me');
        this.session.user = me;
        if (me.must_change_pwd) {
          this.showChangePwd = true;
        }
        await this.loadAll();
        this.connectWS();
      } catch (e) {
        localStorage.removeItem('khan_token');
      localStorage.removeItem('khan_readmap');
      localStorage.removeItem('khan_lang');
      // Keep khan_theme (user preference persists across sessions)
        this.session.token = '';
        this.session.user = null;
      }
    },

    async setup() {
      this.loginBusy = true;
      this.loginError = '';
      try {
        const data = await this.api('/panel/setup', {
          method: 'POST',
          body: JSON.stringify({
            admin_username: this.setupForm.username,
            admin_password: this.setupForm.password,
            company_name: this.setupForm.displayName,
          }),
        });
        if (data.ok) {
          this.loginForm.username = this.setupForm.username;
          this.loginForm.password = this.setupForm.password;
          await this.login();
        }
      } catch (e) {
        this.loginError = e.message;
      } finally {
        this.loginBusy = false;
      }
    },

    async login() {
      this.loginBusy = true;
      this.loginError = '';
      try {
        const data = await this.api('/api/auth/login', {
          method: 'POST',
          body: JSON.stringify(this.loginForm),
        });
        this.session.token = data.token;
        this.session.user = data.user;
        localStorage.setItem('khan_token', data.token);
        this.loginForm.password = '';
        this.setLang(this.uiLang);
        await this.loadAll();
        this.connectWS();
      } catch (e) {
        this.loginError = e.message;
      } finally {
        this.loginBusy = false;
      }
    },

    async logout() {
      try { await this.api('/api/auth/logout', { method: 'POST' }); } catch (e) {}
      if (this.ws) this.ws.close();
      localStorage.removeItem('khan_token');
      this.session = { token: '', user: null };
      this.rooms = [];
      this.messages = {};
      this.unreadMap = {};
      this.readMap = {};
      this.currentRoom = null;
      this.currentRoomId = null;
      document.title = this.docTitle;
    },

    /* ─────────── Data loading ─────────── */
    async loadAll() {
      try {
        const [users, rooms] = await Promise.all([
          this.api('/api/users'),
          this.api('/api/rooms'),
        ]);
        this.users = Array.isArray(users) ? users : [];
        this.rooms = Array.isArray(rooms) ? rooms : [];
        try {
          this.readMap = JSON.parse(localStorage.getItem('khan_readmap') || '{}');
        } catch (e) { this.readMap = {}; }

        for (const room of this.rooms) {
          this.loadMembers(room.id).catch(() => {});
          this.loadMessages(room.id).catch(() => {});
        }
        // Load public/private rooms + departments
        this.loadRoomsDiscovery();
        this.loadDepartments();
        this.loadPinned();
      } catch (e) {
        console.error('loadAll error:', e);
      }
    },

    async loadRoomsDiscovery() {
      try {
        const [pub, priv] = await Promise.all([
          this.api('/api/rooms/public'),
          this.api('/api/rooms/private'),
        ]);
        this.publicRooms = Array.isArray(pub) ? pub : [];
        this.privateRooms = Array.isArray(priv) ? priv : [];
      } catch (e) {}
    },

    async loadDepartments() {
      try {
        const depts = await this.api('/api/departments');
        this.departments = Array.isArray(depts) ? depts : [];
      } catch (e) {}
    },

    async loadMembers(roomId) {
      try {
        const members = await this.api(`/api/rooms/${roomId}/members`);
        this.roomMembers[roomId] = Array.isArray(members) ? members : [];
      } catch (e) {}
    },

    async loadMessages(roomId) {
      try {
        const msgs = await this.api(`/api/messages/${roomId}?limit=50`);
        this.messages[roomId] = Array.isArray(msgs) ? msgs : [];
        this.hasMore[roomId] = (msgs && msgs.length) ? msgs.length >= 50 : false;
        this.computeUnread(roomId);
      } catch (e) {}
    },

    async loadPinned() {
      if (!this.currentRoom) return;
      try {
        const pins = await this.api(`/api/messages/${this.currentRoom.id}/pins`);
        this.pinnedMessages = Array.isArray(pins) ? pins : [];
      } catch (e) { this.pinnedMessages = []; }
    },

    computeUnread(roomId) {
      if (roomId === this.currentRoomId) {
        this.unreadMap[roomId] = 0;
        return;
      }
      const msgs = this.messages[roomId] || [];
      const lastRead = this.readMap[roomId] || 0;
      let count = 0;
      for (const m of msgs) {
        if (m.id > lastRead) count++;
      }
      this.unreadMap[roomId] = count;
    },

    async loadOlder() {
      const roomId = this.currentRoomId;
      if (!roomId || this.loadingMore || !this.hasMore[roomId]) return;
      this.loadingMore = true;
      try {
        const msgs = this.messages[roomId] || [];
        const before = msgs.length ? msgs[0].id : 0;
        const older = await this.api(`/api/messages/${roomId}?limit=50&before=${before}`);
        if (older.length) {
          this.messages[roomId] = [...older, ...msgs];
        }
        this.hasMore[roomId] = older.length >= 50;
      } catch (e) {}
      this.loadingMore = false;
    },

    /* ─────────── Rooms ─────────── */
    chatName(room) {
      if (room.type === 'dm') {
        const other = this.dmPartner(room);
        return other ? (other.display_name || other.username) : 'گفتگو';
      }
      return room.name || 'گروه';
    },

    roomTitle(room) {
      return this.chatName(room);
    },

    roomName(roomId) {
      const room = this.rooms.find(r => r.id === roomId);
      return room ? this.chatName(room) : 'اتاق ' + roomId;
    },

    roomIcon(room) {
      if (!room) return '🏠';
      return { dm: '💬', group: '👥', public: '🌐', private: '🔒', channel: '📢' }[room.type] || '🏠';
    },

    dmPartner(room) {
      if (!room || room.type !== 'dm' || !room.id) return null;
      const members = this.roomMembers[room.id] || [];
      const me = this.session.user ? this.session.user.id : -1;
      return members.find(m => m.id !== me) || null;
    },

    isDM(room) { return room && room.type === 'dm'; },
    isUserOnline(user) { return !!(user && user.online); },

    filteredRooms() {
      const q = this.chatSearch.toLowerCase();
      return this.rooms.filter(r => !q || this.chatName(r).toLowerCase().includes(q));
    },

    filteredRoomsByDept(deptId) {
      return this.filteredRooms().filter(r => (r.department_id || 0) === deptId);
    },

    filteredRoomsNoDept() {
      return this.filteredRooms().filter(r => !r.department_id);
    },

    toggleDept(id) {
      this.deptCollapsed[id] = !this.deptCollapsed[id];
    },

    deptRoomCount(id) {
      return this.rooms.filter(r => (r.department_id || 0) === id).length;
    },

    async openRoom(room) {
      this.currentRoom = room;
      this.currentRoomId = room.id;
      this.mobileSidebar = false;
      if (window.matchMedia('(max-width: 900px)').matches) this.mobileChatOpen = true;
      this.showRoomInfo = false;
      const msgs = this.messages[room.id] || [];
      const lastId = msgs.length ? msgs[msgs.length - 1].id : 0;
      this.readMap[room.id] = lastId;
      this.unreadMap[room.id] = 0;
      try { localStorage.setItem('khan_readmap', JSON.stringify(this.readMap)); } catch (e) {}
      this.updateDocTitle();
      await this.loadMessages(room.id);
      await this.loadMembers(room.id);
      await this.loadPinned();
      this.checkUrgent();
      this.$nextTick(() => this.scrollToBottom());
      // mark read via WS
      if (lastId) this.wsSend({ type: 'mark_read', room_id: room.id, last_id: lastId });
    },

    async joinRoom(room) {
      try {
        await this.api(`/api/rooms/${room.id}/join`, { method: 'POST' });
        this.toast('✅ ' + this.t('clickToJoin'), 'success');
        await this.loadAll();
        const fresh = this.rooms.find(r => r.id === room.id);
        if (fresh) await this.openRoom(fresh);
      } catch (e) {
        this.toast(e.message, 'error');
      }
    },

    openUserProfile(u) {
      this.selectedUser = u;
      this.showUserProfile = true;
    },

    async searchUsers() {
      const q = (this.userSearchQuery || '').trim();
      if (!q) return;
      try {
        const res = await this.api('/api/users/search?q=' + encodeURIComponent(q));
        this.userSearchResults = Array.isArray(res) ? res : [];
      } catch (e) {
        this.userSearchResults = [];
      }
    },

    async startDM(u) {
      if (!u) return;
      try {
        const data = await this.api(`/api/rooms/dm/${u.id}`, { method: 'POST' });
        const fresh = this.rooms ? this.rooms.find(r => r.id === data.room_id) : null;
        if (fresh) await this.openRoom(fresh);
        else {
          // reload rooms so the new DM appears, then open it
          await this.loadRooms();
          const again = this.rooms.find(r => r.id === data.room_id);
          if (again) await this.openRoom(again);
        }
      } catch (e) {
        this.toast(e.message, 'error');
      }
    },

    async createRoom() {
      if (!this.createRoomForm.name.trim()) { this.toast(this.t('roomName') + ' لازم است', 'error'); return; }
      try {
        const room = await this.api('/api/rooms', {
          method: 'POST',
          body: JSON.stringify(this.createRoomForm),
        });
        this.rooms.unshift(room);
        this.showCreateRoom = false;
        this.createRoomForm = { name: '', type: 'group', department: 0 };
        this.toast('✅ ' + this.t('createRoom'), 'success');
        await this.openRoom(room);
        this.loadRoomsDiscovery();
      } catch (e) {
        this.toast(e.message, 'error');
      }
    },

    async inviteMember() {
      if (!this.inviteUserId) return;
      try {
        await this.api(`/api/rooms/${this.currentRoom.id}/invite`, {
          method: 'POST',
          body: JSON.stringify({ user_id: this.inviteUserId }),
        });
        this.toast('✅ ' + this.t('addMember'), 'success');
        this.inviteUserId = null;
        await this.loadMembers(this.currentRoom.id);
      } catch (e) {
        this.toast(e.message, 'error');
      }
    },

    async removeMember(m) {
      if (!confirm(this.t('removeMember') + '؟')) return;
      try {
        await this.api(`/api/rooms/${this.currentRoom.id}/members/${m.id}`, { method: 'DELETE' });
        await this.loadMembers(this.currentRoom.id);
      } catch (e) {
        this.toast(e.message, 'error');
      }
    },

    /* ─────────── Messages ─────────── */
    sendMessage() {
      const text = this.draft.trim();
      if (!text) return;
      if (!this.currentRoom) return;

      const mentions = this.extractMentions(text);
      const payload = {
        type: 'send_message',
        room_id: this.currentRoom.id,
        text,
        mentions,
      };
      if (this.replyTarget) {
        payload.reply_to = this.replyTarget.id;
      }
      this.wsSend(payload);
      this.draft = '';
      this.replyTarget = null;
      this.autoGrow();
    },

    sendSticker(s) {
      if (!this.currentRoom) return;
      this.wsSend({ type: 'send_message', room_id: this.currentRoom.id, text: 'sticker:' + s });
      this.showStickers = false;
    },

    extractMentions(text) {
      const mentions = [];
      const re = /@([\w\u0600-\u06FF]+)/g;
      let m;
      while ((m = re.exec(text))) {
        const uname = m[1];
        const u = this.users.find(x => x.username === uname);
        if (u) mentions.push(u.id);
      }
      return mentions;
    },

    insertEmoji(e) {
      this.draft += e;
      this.autoGrow();
      this.showEmojis = false;
    },

    onComposerKeydown(e) {
      // Mention autocomplete navigation
      if (this.mentionPopupVisible) {
        if (e.key === 'ArrowDown') {
          e.preventDefault();
          const idx = this.mentionCandidates.findIndex(c => c.id === this.mentionSelected);
          this.mentionSelected = this.mentionCandidates[(idx + 1) % this.mentionCandidates.length].id;
          return;
        }
        if (e.key === 'ArrowUp') {
          e.preventDefault();
          const idx = this.mentionCandidates.findIndex(c => c.id === this.mentionSelected);
          this.mentionSelected = this.mentionCandidates[(idx - 1 + this.mentionCandidates.length) % this.mentionCandidates.length].id;
          return;
        }
        if (e.key === 'Enter' || e.key === 'Tab') {
          e.preventDefault();
          const sel = this.mentionCandidates.find(c => c.id === this.mentionSelected);
          if (sel) this.insertMention(sel);
          return;
        }
        if (e.key === 'Escape') {
          this.mentionPopupVisible = false;
          return;
        }
      }
      if (e.key === 'Escape') {
        this.replyTarget = null;
        this.editingMsg = null;
        this.showStickers = false;
        this.showEmojis = false;
      }
    },

    onComposerInput() {
      // Detect @mention typing
      const text = this.draft;
      const atIdx = text.lastIndexOf('@');
      if (atIdx >= 0 && atIdx > text.lastIndexOf(' ')) {
        const q = text.slice(atIdx + 1).toLowerCase();
        this.mentionQuery = q;
        this.mentionCandidates = this.currentRoomMembers.filter(u =>
          u.id !== this.session.user.id &&
          (u.username || '').toLowerCase().includes(q) ||
          (u.display_name || '').toLowerCase().includes(q)
        ).slice(0, 8);
        if (this.mentionCandidates.length) {
          this.mentionPopupVisible = true;
          this.mentionSelected = this.mentionCandidates[0].id;
        } else {
          this.mentionPopupVisible = false;
        }
      } else {
        this.mentionPopupVisible = false;
      }
      this.autoGrow();
      // typing indicator (throttled)
      if (this.currentRoom && Date.now() - (this._lastTyping || 0) > 3000) {
        this._lastTyping = Date.now();
        this.wsSend({ type: 'typing', room_id: this.currentRoom.id });
      }
    },

    insertMention(u) {
      const text = this.draft;
      const atIdx = text.lastIndexOf('@');
      if (atIdx >= 0) {
        this.draft = text.slice(0, atIdx) + '@' + u.username + ' ';
      }
      this.mentionPopupVisible = false;
      this.autoGrow();
      this.$refs.composer && this.$refs.composer.focus();
    },

    autoGrow() {
      this.$nextTick(() => {
        const ta = this.$refs.composer;
        if (ta) {
          ta.style.height = 'auto';
          ta.style.height = Math.min(ta.scrollHeight, 120) + 'px';
        }
      });
    },

    startEdit(msg) {
      this.editingMsg = msg;
      this.draft = msg.text;
      this.contextMenu = null;
      this.$nextTick(() => {
        const ta = this.$refs.composer;
        if (ta) ta.focus();
      });
    },

    cancelEdit() {
      this.editingMsg = null;
      this.draft = '';
    },

    saveEdit() {
      if (!this.editingMsg) return;
      this.wsSend({ type: 'edit_message', message_id: this.editingMsg.id, text: this.draft.trim() });
      this.editingMsg = null;
      this.draft = '';
    },

    async deleteMessage() {
      const msg = this.contextMenu ? this.contextMenu.msg : null;
      if (!msg || !confirm(this.t('delete') + '؟')) return;
      this.contextMenu = null;
      try {
        await this.api(`/api/messages/${msg.id}`, { method: 'DELETE' });
        const roomId = this.currentRoomId;
        this.messages[roomId] = (this.messages[roomId] || []).filter(m => m.id !== msg.id);
      } catch (e) {
        this.toast(e.message, 'error');
      }
    },

    canDeleteMessage(msg) {
      if (!msg) return false;
      if (msg.sender_id === this.session.user.id) return true;
      return this.canManageRoom;
    },

    /* Reactions */
    aggregateReactions(msg) {
      const agg = {};
      for (const r of (msg.reactions || [])) {
        agg[r.emoji] = (agg[r.emoji] || 0) + 1;
      }
      return agg;
    },

    reactedByMe(msg, emoji) {
      return (msg.reactions || []).some(r => r.emoji === emoji && r.user_id === this.session.user.id);
    },

    toggleReaction(msg, emoji) {
      const reacted = this.reactedByMe(msg, emoji);
      this.wsSend({
        type: reacted ? 'remove_reaction' : 'add_reaction',
        message_id: msg.id,
        emoji,
      });
    },

    quickReact(msg, emoji) {
      this.wsSend({ type: 'add_reaction', message_id: msg.id, emoji });
      this.contextMenu = null;
    },

    /* ─────────── Context Menu ─────────── */
    openContextMenu(event, msg) {
      this.contextMenu = { msg, x: event.clientX, y: event.clientY };
      event.preventDefault();
    },

    replyToMessage() {
      if (!this.contextMenu) return;
      this.replyTarget = this.contextMenu.msg;
      this.contextMenu = null;
      this.$nextTick(() => {
        const ta = this.$refs.composer;
        if (ta) ta.focus();
      });
    },

    forwardMessage() {
      this.forwardTarget = this.contextMenu ? this.contextMenu.msg : null;
      this.contextMenu = null;
      this.showForwardModal = true;
    },

    async forwardMessageTo(chat) {
      if (!this.forwardTarget) return;
      try {
        await this.api(`/api/messages/${this.forwardTarget.id}/forward`, {
          method: 'POST',
          body: JSON.stringify({ room_id: chat.id }),
        });
        this.toast('✅ ' + this.t('forward'), 'success');
        this.showForwardModal = false;
        if (this.currentRoomId !== chat.id) {
          await this.loadMessages(chat.id);
        }
      } catch (e) {
        this.toast(e.message, 'error');
      }
    },

    async pinMessage() {
      const msg = this.contextMenu ? this.contextMenu.msg : null;
      if (!msg) return;
      this.contextMenu = null;
      this.wsSend({ type: 'pin_message', message_id: msg.id });
      this.toast('📌 ' + this.t('pin'), 'success');
    },

    async unpinMessage() {
      const msg = this.contextMenu ? this.contextMenu.msg : null;
      if (!msg) return;
      this.contextMenu = null;
      this.wsSend({ type: 'unpin_message', message_id: msg.id });
    },

    async toggleUrgent() {
      const msg = this.contextMenu ? this.contextMenu.msg : null;
      if (!msg) return;
      this.contextMenu = null;
      await this.api(`/api/messages/${msg.id}/urgent`, { method: 'POST' });
      msg.urgent = !msg.urgent;
      this.checkUrgent();
    },

    editMessage() {
      const msg = this.contextMenu ? this.contextMenu.msg : null;
      if (!msg) return;
      this.startEdit(msg);
    },

    /* ─────────── Search ─────────── */
    openSearch() {
      this.showSearchOverlay = true;
      this.searchQuery = '';
      this.searchResults = [];
      this.$nextTick(() => {
        const inp = this.showSearchOverlay ? document.querySelector('.search-results-header input') : null;
        if (inp) inp.focus();
      });
    },

    async doSearch() {
      const q = this.searchQuery.trim() || this.smartSearchQuery.trim();
      if (!q) return;
      try {
        const results = await this.api(`/api/messages/search?q=${encodeURIComponent(q)}`);
        this.searchResults = Array.isArray(results) ? results : [];
        if (this.sidebarTab === 'search') {
          this.showSearchOverlay = false;
        }
      } catch (e) {
        this.toast(e.message, 'error');
      }
    },

    gotoMessage(r) {
      const room = this.rooms.find(x => x.id === r.room_id);
      if (room) {
        this.showSearchOverlay = false;
        this.openRoom(room).then(() => {
          this.gotoMessageId(r.id);
        });
      }
    },

    gotoMessageId(msgId) {
      const roomId = this.currentRoomId;
      const msgs = this.messages[roomId] || [];
      let idx = msgs.findIndex(m => m.id === msgId);
      if (idx >= 0) {
        this.$nextTick(() => {
          const area = this.$refs.messagesArea;
          if (!area) return;
          const els = area.querySelectorAll('.msg-group');
          if (els[idx]) {
            els[idx].scrollIntoView({ behavior: 'smooth', block: 'center' });
            els[idx].style.animation = 'flashHighlight 1s';
            setTimeout(() => { els[idx].style.animation = ''; }, 1000);
          }
        });
      } else {
        // Need to load older messages around this id
        this.loadMessages(roomId).then(() => this.gotoMessageId(msgId));
      }
    },

    gotoPinned(pin) {
      this.gotoMessageId(pin.id);
    },

    replyPreview(m) {
      if (!m.reply_text) return m.reply_to || '';
      return m.reply_text.slice(0, 60);
    },

    /* ─────────── Polls ─────────── */
    openPollModal() {
      this.pollForm = { question: '', options: ['', ''] };
      this.showPollModal = true;
    },

    addPollOption() {
      if (this.pollForm.options.length < 10) this.pollForm.options.push('');
    },

    removePollOption(i) {
      if (this.pollForm.options.length > 2) this.pollForm.options.splice(i, 1);
    },

    async createPoll() {
      if (!this.pollForm.question.trim()) { this.toast('سوال لازم است', 'error'); return; }
      const options = this.pollForm.options.map(o => o.trim()).filter(Boolean);
      if (options.length < 2) { this.toast('حداقل ۲ گزینه لازم است', 'error'); return; }
      this.wsSend({
        type: 'poll_create',
        room_id: this.currentRoom.id,
        poll_question: this.pollForm.question.trim(),
        poll_options: options,
      });
      this.showPollModal = false;
    },

    async votePoll(poll, optionIdx) {
      if (poll.closed) return;
      if (this.userVoted(poll, this.session.user.id) === optionIdx) return;
      this.wsSend({ type: 'poll_vote', poll_id: poll.id, option: optionIdx });
    },

    userVoted(poll, uid) {
      if (!poll.votes) return -1;
      for (const opt in poll.votes) {
        if (poll.votes[opt].includes(uid)) return parseInt(opt);
      }
      return -1;
    },

    pollPct(poll, optIdx) {
      const total = this.pollTotalVotes(poll);
      if (!total) return 0;
      const n = (poll.votes && poll.votes[optIdx] ? poll.votes[optIdx].length : 0);
      return Math.round((n / total) * 100);
    },

    pollTotalVotes(poll) {
      let t = 0;
      if (poll.votes) {
        for (const k in poll.votes) t += poll.votes[k].length;
      }
      return t;
    },

    /* ─────────── File upload ─────────── */
    triggerFileUpload() {
      this.$refs.fileInput && this.$refs.fileInput.click();
    },

    async onFileSelect(e) {
      const file = e.target.files[0];
      if (!file) return;
      await this.uploadFile(file);
      e.target.value = '';
    },

    async onFileDrop(e) {
      this.dragOver = false;
      const file = e.dataTransfer.files[0];
      if (!file) return;
      await this.uploadFile(file);
    },

    async uploadFile(file) {
      if (!this.currentRoom) return;
      const form = new FormData();
      form.append('file', file);
      this.uploadingFile = true;
      try {
        const res = await fetch(`/api/rooms/${this.currentRoom.id}/files`, {
          method: 'POST',
          headers: { 'Authorization': 'Bearer ' + this.session.token },
          body: form,
        });
        const data = await res.json();
        if (!res.ok) throw new Error(data.error || 'خطا');
        this.wsSend({
          type: 'send_message',
          room_id: this.currentRoom.id,
          file_id: data.file_id,
          text: '',
        });
        this.toast('📎 ' + file.name, 'success');
      } catch (err) {
        this.toast(err.message, 'error');
      } finally {
        this.uploadingFile = false;
      }
    },

    async uploadAvatar(e) {
      const file = e.target.files[0];
      if (!file) return;
      const form = new FormData();
      form.append('avatar', file);
      try {
        const res = await fetch('/api/users/avatar', {
          method: 'POST',
          headers: { 'Authorization': 'Bearer ' + this.session.token },
          body: form,
        });
        const data = await res.json();
        if (res.ok && data.user) {
          this.session.user = data.user;
          const u = this.users.find(x => x.id === this.session.user.id);
          if (u) u.avatar = data.user.avatar;
          this.toast('✅ ' + this.t('avatarUpload'), 'success');
        } else {
          throw new Error(data.error || 'خطا');
        }
      } catch (err) {
        this.toast(err.message, 'error');
      }
    },

    /* ─────────── Read receipts ─────────── */
    readCount(msg) {
      // Count members who have read this message (approximation: all members who are online)
      return this.currentRoomMembers.filter(m => m.id !== msg.sender_id).length;
    },

    readReceiptTitle(msg) {
      const total = this.currentRoomMembers.length - 1;
      const read = this.readCount(msg);
      return read >= total ? this.t('seen') : this.t('delivered');
    },

    /* ─────────── Urgent ─────────── */
    checkUrgent() {
      const msgs = this.messages[this.currentRoomId] || [];
      this.currentRoomUrgent = msgs.some(m => m.urgent && !m.read_by_me);
    },

    /* ─────────── Text rendering ─────────── */
    renderText(text) {
      if (!text) return '';
      // Escape HTML
      let t = text.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
      // Stickers (sent as sticker:EMOJI)
      if (t.startsWith('sticker:')) {
        const emoji = t.replace('sticker:', '').trim();
        return `<span style="font-size:56px; line-height:1;">${emoji}</span>`;
      }
      // Escape quotes to prevent attribute injection (XSS)
      t = t.replace(/"/g, '&quot;').replace(/'/g, '&#39;');
      // Links (exclude quotes from URL chars to prevent XSS)
      t = t.replace(/(https?:\/\/[^\s<"']+)/g, '<a href="$1" target="_blank" rel="noopener" style="color:var(--accent-hover); text-decoration:underline;">$1</a>');
      // Mentions
      t = t.replace(/@([\w\u0600-\u06FF]+)/g, '<span class="msg-mention">@$1</span>');
      // Code
      t = t.replace(/`([^`]+)`/g, '<code style="background:var(--bg-elevated); padding:1px 6px; border-radius:4px; font-size:0.9em; direction:ltr; display:inline-block;">$1</code>');
      // Bold
      t = t.replace(/\*\*([^*]+)\*\*/g, '<b>$1</b>');
      return t;
    },

    /* ─────────── Time (Persian) ─────────── */
    dayKey(dateStr) {
      if (!dateStr) return '';
      const d = new Date(dateStr);
      const now = new Date();
      const yesterday = new Date(now);
      yesterday.setDate(now.getDate() - 1);
      if (d.toDateString() === now.toDateString()) return 'today';
      if (d.toDateString() === yesterday.toDateString()) return 'yesterday';
      return this.shamsiKey(d);
    },

    dayLabel(key) {
      if (key === 'today') return this.t('today');
      if (key === 'yesterday') return this.t('yesterday');
      return key;
    },

    shamsiKey(d) {
      try {
        return new Intl.DateTimeFormat('fa-IR', { weekday: 'long', day: 'numeric', month: 'long' }).format(d);
      } catch (e) {
        return d.toLocaleDateString('fa-IR');
      }
    },

    timeShort(dateStr) {
      if (!dateStr) return '';
      try {
        return new Intl.DateTimeFormat('fa-IR', { hour: '2-digit', minute: '2-digit' }).format(new Date(dateStr));
      } catch (e) {
        return new Date(dateStr).toLocaleTimeString('fa-IR', { hour: '2-digit', minute: '2-digit' });
      }
    },

    fullTime(dateStr) {
      if (!dateStr) return '';
      try {
        return new Intl.DateTimeFormat('fa-IR', { year: 'numeric', month: 'long', day: 'numeric', hour: '2-digit', minute: '2-digit' }).format(new Date(dateStr));
      } catch (e) {
        return dateStr;
      }
    },

    relativeTime(dateStr) {
      if (!dateStr) return '';
      const d = new Date(dateStr);
      const now = Date.now();
      const diff = now - d.getTime();
      const min = Math.floor(diff / 60000);
      const hour = Math.floor(diff / 3600000);
      const day = Math.floor(diff / 86400000);
      if (min < 1) return this.t('justNow');
      if (min < 60) return min + ' ' + this.t('minAgo');
      if (hour < 24) return hour + ' ' + this.t('hourAgo');
      if (day < 7) return day + ' ' + this.t('dayAgo');
      return this.fullTime(dateStr).split('،')[0];
    },

    /* ─────────── Avatar ─────────── */
    avatarStyle(user) {
      if (!user) return { background: '#4d9de0' };
      if (user.avatar) {
        return { backgroundImage: `url(${user.avatar})`, backgroundSize: 'cover', backgroundPosition: 'center' };
      }
      const colors = ['#e17055', '#0984e3', '#00b894', '#6c5ce7', '#fdcb6e', '#e84393', '#00cec9', '#d63031'];
      const name = (user.display_name || user.username || '?');
      let hash = 0;
      for (const ch of name) hash = (hash * 31 + ch.charCodeAt(0)) % 997;
      return { background: colors[hash % colors.length] };
    },

    initials(name) {
      if (!name) return '?';
      const parts = String(name).trim().split(/\s+/);
      if (parts.length >= 2) return (parts[0][0] + parts[1][0]).toUpperCase();
      return name[0].toUpperCase();
    },

    safeUser(user) {
      return user || {};
    },

    roleLabel(role) {
      return { user: 'کاربر', sup: 'سوپروایزر', adm: 'ادمین', sadm: 'مدیر اصلی' }[role] || role;
    },

    memberRoleLabel(userId) {
      if (!this.currentRoom) return '';
      const role = (this.currentRoom.member_roles && this.currentRoom.member_roles[userId]) || '';
      return { owner: '👑 سازنده', admin: '🛡 ادمین', member: '👤' }[role] || '';
    },

    /* ─────────── Scroll ─────────── */
    scrollToBottom() {
      this.$nextTick(() => {
        const area = this.$refs.messagesArea;
        if (area) area.scrollTop = area.scrollHeight;
      });
    },

    onMessagesScroll() {
      const area = this.$refs.messagesArea;
      if (area && area.scrollTop < 60) {
        this.loadOlder();
      }
    },

    /* ─────────── Toasts ─────────── */
    toast(text, type = 'info') {
      const id = ++this.toastId;
      this.toasts.push({ id, text, type });
      setTimeout(() => this.removeToast(id), 3500);
    },

    removeToast(id) {
      this.toasts = this.toasts.filter(t => t.id !== id);
    },

    /* ─────────── Profile ─────────── */
    async saveProfile() {
      try {
        const me = await this.api('/api/users/me', {
          method: 'PUT',
          body: JSON.stringify({ display_name: this.profileForm.display_name }),
        });
        this.session.user = me;
        const u = this.users.find(x => x.id === me.id);
        if (u) u.display_name = me.display_name;
        this.showProfile = false;
        this.toast('✅ ' + this.t('save'), 'success');
      } catch (e) {
        this.toast(e.message, 'error');
      }
    },

    async changePassword() {
      try {
        await this.api('/api/auth/change-password', {
          method: 'POST',
          body: JSON.stringify({
            current_password: this.pwdForm.current,
            new_password: this.pwdForm.next,
          }),
        });
        this.showChangePwd = false;
        this.pwdForm = { current: '', next: '' };
        if (this.session.user) this.session.user.must_change_pwd = false;
        this.toast('✅ ' + this.t('changePassword'), 'success');
      } catch (e) {
        this.toast(e.message, 'error');
      }
    },

    /* ─────────── Admin: Users ─────────── */
    canPromote(u) {
      if (!u || !this.session.user) return false;
      return this.isAdmin && u.role !== 'adm' && u.role !== 'sadm' && u.id !== this.session.user.id;
    },

    async createUser() {
      try {
        const u = await this.api('/api/users', {
          method: 'POST',
          body: JSON.stringify(this.createUserForm),
        });
        this.users.push(u);
        this.showCreateUser = false;
        this.createUserForm = { username: '', display_name: '', password: '', role: 'user' };
        this.toast('✅ ' + this.t('newUser'), 'success');
      } catch (e) {
        this.toast(e.message, 'error');
      }
    },

    async toggleRole(u) {
      const newRole = u.role === 'sup' ? 'user' : 'sup';
      try {
        await this.api(`/api/users/${u.id}/role`, {
          method: 'POST',
          body: JSON.stringify({ role: newRole }),
        });
        u.role = newRole;
      } catch (e) {
        this.toast(e.message, 'error');
      }
    },

    async resetPassword(u) {
      const pwd = prompt('رمز جدید برای ' + (u.display_name || u.username) + ':');
      if (!pwd) return;
      try {
        await this.api(`/api/users/${u.id}/reset-password`, {
          method: 'POST',
          body: JSON.stringify({ new_password: pwd }),
        });
        this.toast('✅ ' + this.t('changePassword'), 'success');
      } catch (e) {
        this.toast(e.message, 'error');
      }
    },

    async toggleActive(u) {
      try {
        await this.api(`/api/users/${u.id}/toggle-active`, { method: 'POST' });
        u.active = !u.active;
      } catch (e) {
        this.toast(e.message, 'error');
      }
    },

    async deleteUser(u) {
      if (!confirm(this.t('delete') + ' ' + (u.display_name || u.username) + '؟')) return;
      try {
        await this.api(`/api/users/${u.id}`, { method: 'DELETE' });
        this.users = this.users.filter(x => x.id !== u.id);
        this.toast('🗑 ' + this.t('delete'), 'success');
      } catch (e) {
        this.toast(e.message, 'error');
      }
    },

    /* ─────────── Admin: Departments ─────────── */
    async createDepartment() {
      if (!this.createDeptForm.name.trim()) return;
      try {
        const d = await this.api('/api/departments', {
          method: 'POST',
          body: JSON.stringify({ name: this.createDeptForm.name.trim() }),
        });
        this.departments.push(d);
        this.showCreateDept = false;
        this.createDeptForm = { name: '' };
        this.toast('✅ ' + this.t('newDepartment'), 'success');
      } catch (e) {
        this.toast(e.message, 'error');
      }
    },

    async deleteDepartment(d) {
      if (!confirm(this.t('delete') + ' ' + d.name + '؟')) return;
      try {
        await this.api(`/api/departments/${d.id}`, { method: 'DELETE' });
        this.departments = this.departments.filter(x => x.id !== d.id);
      } catch (e) {
        this.toast(e.message, 'error');
      }
    },

    /* ─────────── Admin: License ─────────── */
    async loadLicense() {
      try {
        const lic = await this.api('/api/settings/license');
        this.license = lic;
      } catch (e) {}
    },

    async applyLicense(event) {
      const file = event.target.files[0];
      if (!file) return;
      const form = new FormData();
      form.append('license', file);
      try {
        const res = await fetch('/api/settings/license', {
          method: 'POST',
          headers: { 'Authorization': 'Bearer ' + this.session.token },
          body: form,
        });
        const data = await res.json();
        this.license = data;
        if (data.applied) {
          this.toast('✅ ' + this.t('applyLicense') + ' — ' + this.t('maxUsers') + ': ' + data.max_users, 'success');
        } else {
          this.toast('⚠️ ' + (data.error || ''), 'error');
        }
      } catch (e) {
        this.toast('خطا', 'error');
      }
    },

    async removeLicense() {
      if (!confirm(this.t('removeLicense') + '؟')) return;
      try {
        await this.api('/api/settings/license', { method: 'DELETE' });
        this.license = { state: 'free', max_users: 20, licensed_to: '', expiry: '', error: '' };
        this.toast('✅', 'success');
      } catch (e) {
        this.toast(e.message, 'error');
      }
    },

    /* ─────────── Admin: Network ─────────── */
    async loadNetwork() {
      try {
        const n = await this.api('/api/settings/network');
        this.networkForm = { address_type: n.address_type, ip: n.ip, dns: n.dns, port: n.port };
        this.networkInfo = n;
      } catch (e) {}
    },

    async saveNetwork() {
      try {
        await this.api('/api/settings/network', {
          method: 'POST',
          body: JSON.stringify(this.networkForm),
        });
        this.toast('✅ ' + this.t('saveSettings'), 'success');
      } catch (e) {
        this.toast(e.message, 'error');
      }
    },

    /* ─────────── Admin: Backup & Logs ─────────── */
    async manualBackup() {
      try {
        const data = await this.api('/api/settings/backup', { method: 'POST' });
        this.toast('✅ ' + (data.file || this.t('backupNow')), 'success');
        this.loadBackups();
      } catch (e) {
        this.toast(e.message, 'error');
      }
    },

    async loadBackups() {
      try {
        const b = await this.api('/api/settings/backups');
        this.backups = Array.isArray(b) ? b : [];
      } catch (e) {}
    },

    async downloadBackup(name) {
      try {
        const res = await fetch('/api/settings/backups/' + encodeURIComponent(name) + '/download', {
          headers: { 'Authorization': 'Bearer ' + this.session.token },
        });
        const blob = await res.blob();
        const a = document.createElement('a');
        a.href = URL.createObjectURL(blob);
        a.download = name;
        a.click();
        URL.revokeObjectURL(a.href);
      } catch (e) {}
    },

    async restoreBackup(name) {
      if (!confirm(this.t('restoreConfirm'))) return;
      try {
        const data = await this.api('/api/settings/restore', {
          method: 'POST',
          body: JSON.stringify({ name }),
        });
        this.toast('\u2714 ' + (data.restored || this.t('restored')), 'success');
        this.loadBackups();
      } catch (e) {
        this.toast(e.message, 'error');
      }
    },

    async exportUsers() {
      try {
        const res = await fetch('/api/users/export', {
          headers: { 'Authorization': 'Bearer ' + this.session.token },
        });
        const blob = await res.blob();
        const a = document.createElement('a');
        a.href = URL.createObjectURL(blob);
        a.download = 'khan-users-export.csv';
        a.click();
        URL.revokeObjectURL(a.href);
        this.toast('\u271c ' + this.t('exportDone'), 'success');
      } catch (e) {
        this.toast(e.message, 'error');
      }
    },

    async refreshLogs() {
      try {
        const logs = await this.api('/api/settings/logs');
        this.serverLogs = logs.logs || '';
      } catch (e) {
        this.serverLogs = '';
      }
    },

    /* ─────────── Admin Panel ─────────── */
    openAdminPanel() {
      this.showAdminPanel = true;
      this.adminTab = 'users';
      this.loadLicense();
      this.loadNetwork();
      this.loadBackups();
      this.refreshLogs();
      this.loadDepartments();
    },

    /* ─────────── WebSocket ─────────── */
    connectWS() {
      if (!this.session.token) return;
      if (this.ws) { try { this.ws.close(); } catch (e) {} }
      const proto = location.protocol === 'https:' ? 'wss' : 'ws';
      this.ws = new WebSocket(`${proto}://${location.host}/ws?token=${encodeURIComponent(this.session.token)}`);

      this.ws.onopen = () => {
        this.wsRetry = 0;
        console.log('✅ Khan WS connected');
      };

      this.ws.onmessage = (e) => {
        try {
          const ev = JSON.parse(e.data);
          this.handleEvent(ev);
        } catch (err) {
          console.error('bad ws message', err);
        }
      };

      this.ws.onclose = () => {
        if (this.session.token) {
          clearTimeout(this.reconnectTimer);
          this.reconnectTimer = setTimeout(() => this.connectWS(), Math.min(30000, 1000 * Math.pow(2, this.wsRetry++)));
        }
      };

      this.ws.onerror = (e) => { console.error('ws error', e); };
    },

    wsSend(obj) {
      if (this.ws && this.ws.readyState === WebSocket.OPEN) {
        this.ws.send(JSON.stringify(obj));
      }
    },

    handleEvent(ev) {
      switch (ev.type) {
        case 'message': {
          const roomId = ev.room_id;
          if (!this.messages[roomId]) this.messages[roomId] = [];
          // Avoid duplicate
          const p = ev.payload;
          if (!this.messages[roomId].some(m => m.id === p.id)) {
            this.messages[roomId].push(p);
          }
          if (this.currentRoomId === roomId) {
            this.scrollToBottom();
            this.readMap[roomId] = p.id;
            this.unreadMap[roomId] = 0;
            try { localStorage.setItem('khan_readmap', JSON.stringify(this.readMap)); } catch (e) {}
            this.checkUrgent();
          } else {
            this.unreadMap[roomId] = (this.unreadMap[roomId] || 0) + 1;
            this.notifyNewMessage(p, roomId);
          }
          if (!this.rooms.some(r => r.id === roomId)) {
            this.loadAll();
          } else {
            const room = this.rooms.find(r => r.id === roomId);
            if (room) {
              room.lastMessage = p.text || '📎';
              room.lastMessageAt = p.created_at;
              this.rooms = this.rooms.filter(r => r.id !== roomId);
              this.rooms.unshift(room);
            }
          }
          this.updateDocTitle();
          break;
        }
        case 'message_edited': {
          const roomId = ev.room_id;
          const msgs = this.messages[roomId] || [];
          const idx = msgs.findIndex(m => m.id === ev.payload.id);
          if (idx >= 0) {
            msgs[idx].text = ev.payload.text;
            msgs[idx].edited_at = ev.payload.edited_at;
          }
          break;
        }
        case 'message_deleted': {
          const roomId = ev.room_id;
          this.messages[roomId] = (this.messages[roomId] || []).filter(m => m.id !== ev.payload);
          this.pinnedMessages = this.pinnedMessages.filter(p => p.id !== ev.payload);
          // a pinned/urgent notification for this message should no longer exist
          this.closeRoomNotifications(roomId);
          this.checkUrgent();
          break;
        }
        case 'reaction': {
          const roomId = ev.room_id;
          const msgs = this.messages[roomId] || [];
          const msg = msgs.find(m => m.id === ev.payload.message_id);
          if (!msg) break;
          if (!msg.reactions) msg.reactions = [];
          if (ev.payload.remove) {
            msg.reactions = msg.reactions.filter(r => !(r.user_id === ev.payload.user_id && r.emoji === ev.payload.emoji));
          } else {
            if (!msg.reactions.some(r => r.user_id === ev.payload.user_id && r.emoji === ev.payload.emoji)) {
              msg.reactions.push({ user_id: ev.payload.user_id, emoji: ev.payload.emoji });
            }
          }
          break;
        }
        case 'typing': {
          const roomId = ev.room_id;
          if (roomId !== this.currentRoomId) break;
          const uid = ev.payload.user_id;
          if (uid === this.session.user.id) break;
          // keep the sender's name so the indicator can show WHO is typing
          this.typingUsers[uid] = true;
          this.typingNames_[uid] = (ev.payload.name || '').trim();
          clearTimeout(this.typingUsers[uid + '_t']);
          this.typingUsers[uid + '_t'] = setTimeout(() => {
            delete this.typingUsers[uid];
            delete this.typingNames_[uid];
          }, 2500);
          break;
        }
        case 'presence': {
          const p = ev.payload;
          const u = this.users.find(x => x.id === p.user_id);
          if (u) u.online = p.online;
          const rm = this.currentRoomMembers.find(x => x.id === p.user_id);
          if (rm) rm.online = p.online;
          break;
        }
        case 'read_receipt': {
          // update read status visually (single tick → double tick)
          const p = ev.payload;
          const msgs = this.messages[p.room_id] || [];
          for (const m of msgs) {
            if (m.id <= p.last_id && m.sender_id === this.session.user.id) {
              m.read_by = m.read_by || [];
              if (!m.read_by.includes(p.user_id)) m.read_by.push(p.user_id);
            }
          }
          break;
        }
        case 'pin': {
          const p = ev.payload;
          this.toast('📌 ' + this.t('pin'), 'info');
          this.loadPinned();
          break;
        }
        case 'unpin': {
          const p = ev.payload;
          this.pinnedMessages = this.pinnedMessages.filter(x => x.id !== p.message_id);
          break;
        }
        case 'poll': {
          const roomId = ev.room_id;
          if (!this.messages[roomId]) this.messages[roomId] = [];
          const p = ev.payload;
          if (!this.messages[roomId].some(m => m.id === p.id)) {
            this.messages[roomId].push(p);
          }
          break;
        }
        case 'poll_update': {
          const p = ev.payload;
          for (const roomId in this.messages) {
            for (const m of this.messages[roomId]) {
              if (m.poll && m.poll.id === p.poll_id) {
                if (p.votes) m.poll.votes = p.votes;
                if (p.closed) m.poll.closed = true;
              }
            }
          }
          break;
        }
        case 'force_logout':
          this.toast('خروج اجباری', 'error');
          this.logout();
          break;
        case 'error':
          this.toast(ev.payload.error || 'خطا', 'error');
          break;
      }
    },

    /* ─────────── Notifications & Unread ─────────── */
    updateDocTitle() {
      const total = this.totalUnread;
      document.title = total > 0 ? `(${total}) ${this.docTitle}` : this.docTitle;
    },

    requestNotifPermission() {
      if ('Notification' in window && Notification.permission === 'default') {
        Notification.requestPermission().then(p => {
          this.notifEnabled = p === 'granted';
        });
      } else if ('Notification' in window) {
        this.notifEnabled = Notification.permission === 'granted';
      }
    },

    notifyNewMessage(msg, roomId) {
      // Browser notification (only when tab hidden or permission granted)
      if (msg.urgent || ('Notification' in window && Notification.permission === 'granted' && document.hidden)) {
        try {
          const room = this.rooms.find(r => r.id === roomId);
          const title = msg.urgent ? '🚨 ' + this.chatName(room || { type: 'group' }) : this.chatName(room || { type: 'group' });
          const body = (msg.sender_name ? msg.sender_name + ': ' : '') + (msg.text || '📎 فایل');
          const n = new Notification(title, { body, icon: '/img/khan-logo.jpg', tag: 'khan-' + roomId, requireInteraction: !!msg.urgent });
          n.onclick = () => {
            window.focus();
            if (room) {
              const target = this.messages[room.id] || [];
              // open the room, then jump to the exact message if we have it
              if (this.currentRoom && this.currentRoom.id === room.id && this.messages[this.currentRoom.id]) {
                this.scrollToMessage(msg.id);
              } else {
                this.openRoom(room).then(() => {
                  if (this.messages[room.id]) this.scrollToMessage(msg.id);
                });
              }
            }
            n.close();
          };
        } catch (e) {}
      }
    },

    // Scroll the messages pane so the given message id is in view (used by notif click)
    scrollToMessage(id) {
      this.$nextTick(() => {
        const el = this.$refs.messagesArea;
        if (!el) return;
        const target = el.querySelector(`[data-msg-id="${id}"]`);
        if (target) target.scrollIntoView({ block: 'center', behavior: 'smooth' });
      });
    },

    // Close any live desktop notification for a room (used when its message is deleted)
    closeRoomNotifications(roomId) {
      if (!('Notification' in window)) return;
      try {
        // close all notifications tagged for this room
        for (const n of Notification.getNotifications ? Notification.getNotifications() : []) {
          if (n.tag === 'khan-' + roomId) n.close();
        }
      } catch (e) {}
    },
  },

  async mounted() {
    document.documentElement.setAttribute('data-theme', this.uiTheme);
    this.setLang(this.uiLang);
    this.requestNotifPermission();
    try {
      this.serverInfo = await this.api('/api/settings/info');
    } catch (e) {}
    await this.checkSetup();
  },

  watch: {
    showAdminPanel(v) {
      if (v) {
        this.loadLicense();
        this.loadNetwork();
      }
    },
    showProfile(v) {
      if (v && this.session.user) {
        this.profileForm.display_name = this.session.user.display_name || '';
      }
    },
    currentRoomId() {
      this.pinnedMessages = [];
    },
  },
});

app.mount('#app');

/* ============ KEYBOARD SHORTCUTS ============ */
document.addEventListener('keydown', (e) => {
  // Cmd+K / Ctrl+K: focus search
  if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
    e.preventDefault();
    const searchInput = document.querySelector('.search-input, #search-input, input[type="search"]');
    if (searchInput) searchInput.focus();
  }
  // Escape: close modals/drawers
  if (e.key === 'Escape') {
    const drawer = document.getElementById('member-drawer');
    if (drawer && drawer.classList.contains('open')) {
      closeDrawer();
    }
    // Close any open modal
    document.querySelectorAll('.modal.open').forEach(m => m.classList.remove('open'));
  }
});
