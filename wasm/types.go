package main

import (
	"encoding/json"
	store "github.com/StreamSpace/ss-store"
	"time"
)

type Event struct {
	Result struct {
		Topic string `json:"topic,omitempty"`
		Val   string `json:"val,Out,omitempty"`
	}
}

type Out struct {
	Status  int         `json:"status"`
	Message string      `json:"message"` // Message is for the curators
	Data    interface{} `json:"data,omitempty"`
	Details string      `json:"details,omitempty"` // For Debugging
}

type ID struct {
	PeerID    string   `json:"id,omitempty"`
	Publickey string   `json:"PublicKey,omitempty"`
	Addresses []string `json:"Addresses,omitempty"`
}

type SwarmPeers struct {
	Peers []string
}

type Profile struct {
	Id              string `json:"_id,omitempty"`
	Email           string `json:"email,omitempty"`
	FirstName       string `json:"firstName,omitempty"`
	PhoneNumber     string `json:"phoneNumber,omitempty"`
	Role            string `json:"role,omitempty"`
	MFAType         int64  `json:"mfaType,omitempty"`
	IsMFAEnabled    bool   `json:"isMfaEnabled,omitempty"`
	LastLoginAt     string `json:"lastLoginAt,omitempty"`
	IsEmailVerified bool   `json:"isEmailVerified",omitempty"`
}

type Bandwidth struct {
	Incoming float64 `json:"RateIn,omitempty"`
	Outgoing float64 `json:"RateOut,omitempty"`
	Time     int64   `json:"Time,omitempty"`
}

type Version struct {
	AppVersion    string `json:"appversion,omitempty"`
	CurrentCommit string `json:"currentcommit,omitempty"`
	Debug         string `json:"debug,omitempty"`
	Environment   string `json:"environment,omitempty"`
	Epoch         string `json:"epoch,omitempty"`
	CycleDuration string `json:"cycleduration,omitempty"`
}

type AuthResponse struct {
	User  User  `json:"user"`
	Token Token `json:"token"`
}

func (a *AuthResponse) GetNamespace() string {
	return "Auth"
}
func (a *AuthResponse) GetId() string {
	return "1"
}
func (a *AuthResponse) Marshal() ([]byte, error) {
	return json.Marshal(a)
}
func (a *AuthResponse) Unmarshal(val []byte) error {
	return json.Unmarshal(val, a)
}

type Token struct {
	Token        string `json:"accessToken,omitempty"`
	RefreshToken string `jon:"refreshToken,omitempty"`
	Expired      int64  `json:"expiresIn,omitempty"`
}
type User struct {
	Id            string `json:"_id"`
	Email         string `json:"email"`
	FirstName     string `json:"firstName,omitempty"`
	LastName      string `json:"lastName,omitempty"`
	PhoneNumber   string `json:"phoneNumber,omitempty"`
	Role          string `json:"role,omitempty"`
	MFAType       int    `json:"mfaType,omitempty"`
	MFAEnabled    bool   `json:"isMfaEnabled,omitempty"`
	LastLogin     string `json:"lastLoginAt,omitempty"`
	KYCVerified   bool   `json:"isKycVerified,omitempty"`
	EmailVerified bool   `json:"isEmailVerified,omitempty"`
}

func (u *User) GetNamespace() string {
	return "User"
}
func (u *User) GetId() string {
	return u.Id
}
func (u *User) Marshal() ([]byte, error) {
	return json.Marshal(u)
}
func (u *User) Unmarshal(val []byte) error {
	return json.Unmarshal(val, u)
}

type DeviceOwner struct {
	Email string
}

func (do *DeviceOwner) GetNamespace() string {
	return "DeviceOwner"
}
func (do *DeviceOwner) GetId() string {
	return "1"
}
func (do *DeviceOwner) Marshal() ([]byte, error) {
	return []byte(do.Email), nil
}
func (do *DeviceOwner) Unmarshal(val []byte) error {
	do.Email = string(val)
	return nil
}

type Settlement struct {
	Cycle int64     `json:"bcn"`
	Date  time.Time `json:"settlementDate"`
	Rate  float64   `json:"dataRatePerByte"`
}

func (s *Settlement) GetNamespace() string {
	return "Settlement"
}
func (s *Settlement) GetId() string {
	return "1"
}
func (s *Settlement) Marshal() ([]byte, error) {
	return json.Marshal(s)
}
func (s *Settlement) Unmarshal(val []byte) error {
	return json.Unmarshal(val, s)
}

type Balance struct {
	UserId  string  `json:"userId,omitempty"`
	Balance float64 `json:"balance,omitempty"`
	Message string  `json:"message,omitempty"`
}

func (b *Balance) GetNamespace() string {
	return "UserBalance"
}
func (b *Balance) GetId() string {
	return "1"
}
func (b *Balance) Marshal() ([]byte, error) {
	return json.Marshal(b)
}
func (b *Balance) Unmarshal(val []byte) error {
	return json.Unmarshal(val, b)
}

type BCNBalance struct {
	Owned           float64 `json:"owned"`
	Owe             float64 `json:"owe"`
	BytesServed     int64   `json:"served"`
	BytesDownloaded int64   `json:"downloaded"`
	Id              string  `json:"id"`
}

func (b *BCNBalance) GetNamespace() string {
	return "BCNBalance"
}
func (b *BCNBalance) GetId() string {
	return b.Id
}
func (b *BCNBalance) Marshal() ([]byte, error) {
	return json.Marshal(b)
}
func (b *BCNBalance) Unmarshal(val []byte) error {
	return json.Unmarshal(val, b)
}
func (b *BCNBalance) Factory() store.SerializedItem {
	return &BCNBalance{}
}

type Settings struct {
	NodeIndex                      float64 `json:"nodeIndex,omitempty"`
	DeviceID                       string  `json:"deviceId,omitempty"`
	Name                           string  `json:"name,omitempty"`
	Location                       string  `json:"location,omitempty"`
	IPAddress                      string  `json:"ipAddress,omitempty"`
	MaxStorage                     float64 `json:"maxStorage,omitempty"`
	UsedStorage                    float64 `json:"usedStorage"`
	PinnedStorage                  float64 `json:"pinned_storage"`
	HiveStorage                    float64 `json:"hive_storage"`
	PeerID                         string  `json:"peerId,omitempty"`
	PublicKey                      string  `json:"publicKey,omitempty"`
	IsReachable                    bool    `json:"isReachable"`
	IsDNSEligible                  bool    `json:"isDnsEligible"`
	DesktopApplicationNotification bool    `json:"isOSNotification"`
	DesktopApplicationAutoStart    bool    `json:"isAutoStartEnabled"`
	DNS                            string  `json:"dnsAddress"`
	Role                           string  `json:"role,omitempty"`
}

func (b *Settings) GetNamespace() string {
	return "UserSettings"
}
func (b *Settings) GetId() string {
	return "1"
}
func (b *Settings) Marshal() ([]byte, error) {
	return json.Marshal(b)
}
func (b *Settings) Unmarshal(val []byte) error {
	return json.Unmarshal(val, b)
}

type SwarmURL struct {
	URL string `json:"url"`
}
type Location struct {
	City string `json:"city"`
}
type TaskStatus struct {
	Id               int
	Name             string
	Status           string
	AdditionalStatus string
}
type TaskWithProgressStatus struct {
	Description string
	FileName    string
	Progress    float64
}

type ServerStatus struct {
	Rpc   string `json:"Rpc"`
	Http  string `json:"Http"`
	Proxy string `json:"Proxy"`
}

type Config struct {
	APIPort string `json:"APIPort, omitempty"`
	AutoGC string `json:"AutoGC, omitempty"`
	Bootstraps []string `json:"Bootstraps, omitempty"`
	DNS4 string `json:"DNS4, omitempty"`
	DataStore string `json:"DataStore, omitempty"`
	DesktopApplicationAutoStart bool `json:"DesktopApplicationAutoStart, omitempty"`
	DesktopApplicationNotification bool `json:"DesktopApplicationNotification, omitempty"`
	DeviceName string `json:"DeviceName, omitempty"`
	EnableDynamicDNS bool `json:"EnableDynamicDNS, omitempty"`
	EnableFileShare bool `json:"EnableFileShare, omitempty"`
	EnableHop bool `json:"EnableHop, omitempty"`
	GCPeriod string `json:"GCPeriod, omitempty"`
	GatewayPort string `json:"GatewayPort, omitempty"`
	IP4 string `json:"IP4, omitempty"`
	IP6 string `json:"IP6, omitempty"`
	identity Identity `json:"Identity, omitempty"`
	MaxPeers int `json:"MaxPeers, omitempty"`
	ProxyPort string `json:"ProxyPort, omitempty"`
	ReproviderInterval string `json:"ReproviderInterval, omitempty"`
	Storage int `json:"Storage, omitempty"`
	StorageGCWatermark int `json:"StorageGCWatermark, omitempty"`
	Store string `json:"Store, omitempty"`
	SwarmPort string `json:"SwarmPort, omitempty"`
	WebsocketPort string `json:"WebsocketPort, omitempty"`
}

type Identity struct {
	PeerId string `json:"PeerId, omitempty"`
	PrivKey string `json:"PrivKey, omitempty"`
}

type Status struct {
	LoggedIn              bool
	DaemonRunning         bool
	TotalUptimePercentage UptimePercentage
	SessionStartTime      int64
	TaskManagerStatus     []TaskStatus
	ServerDetails         ServerStatus `json:"ServerStatus"`
}
type UptimePercentage struct {
	Status               bool
	Percentage           float64
	SecondsFromInception int64
	Timestamp            int64
}

func (u *UptimePercentage) GetNamespace() string {
	return "UptimePercentage"
}
func (u *UptimePercentage) GetId() string {
	return "1"
}
func (u *UptimePercentage) Marshal() ([]byte, error) {
	return json.Marshal(u)
}
func (u *UptimePercentage) Unmarshal(val []byte) error {
	return json.Unmarshal(val, u)
}

type SessionStartTime struct {
	Timestamp int64
}

func (u *SessionStartTime) GetNamespace() string {
	return "UptimeThisSession"
}
func (u *SessionStartTime) GetId() string {
	return "1"
}
func (u *SessionStartTime) Marshal() ([]byte, error) {
	return json.Marshal(u)
}
func (u *SessionStartTime) Unmarshal(val []byte) error {
	return json.Unmarshal(val, u)
}

type Device struct {
	Name   string `json:"name"`
	PeerId string `json:"peerId"`
}
type Earning struct {
	Earned   float64 `json:"earned"`
	Served   float64 `json:"served"`
	Download float64 `json:"download"`
}
type NetEarnings struct {
	BillingCycles []string             `json:"billingCycles"`
	Devices       []Device             `json:"devices"`
	Data          map[string][]Earning `json:"earnings"`
	CycleStats    [6] CycleStat
	DeviceTotal   Earning
}

type GraphDetails struct {
	BillingCyclesReverse []string  `json:"billingCyclesReverse"`
	Earnings 			 []float64 `json:"Earnings"`
}

func (n *NetEarnings) GetNamespace() string {
	return "NetEarnings"
}
func (n *NetEarnings) GetId() string {
	return "1"
}
func (n *NetEarnings) Marshal() ([]byte, error) {
	return json.Marshal(n)
}
func (n *NetEarnings) Unmarshal(val []byte) error {
	return json.Unmarshal(val, n)
}

type CycleStat struct {
	Cycle      string
	Earned     float64
	Downloaded float64
	Served     float64
}
type FileObj struct {
	Filename               string `json:"filename"`
	Hash                   string `json:"hash"`
	Size                   int64  `json:"size"`
	ShareableEncodedString string `json:"shareable"`
	CreatedAt              int64  `json:"created_at"`
	UpdatedAt              int64  `json:"updated_at"`
	Shared                 bool   `json:"shared"`
	IsPinned               bool   `json:"is_pinned"`
}

func (f *FileObj) GetNamespace() string {
	return "FileObj"
}
func (f *FileObj) GetId() string {
	return f.Filename
}
func (f *FileObj) Marshal() ([]byte, error) {
	return json.Marshal(f)
}
func (f *FileObj) Unmarshal(val []byte) error {
	return json.Unmarshal(val, f)
}
func (f *FileObj) Factory() store.SerializedItem {
	return &FileObj{}
}
func (f *FileObj) SetCreated(unixTime int64) { f.CreatedAt = unixTime }
func (f *FileObj) SetUpdated(unixTime int64) { f.UpdatedAt = unixTime }
func (f *FileObj) GetCreated() int64         { return f.CreatedAt }
func (f *FileObj) GetUpdated() int64         { return f.UpdatedAt }

type FileStatus string

var (
	Selected FileStatus = "selected"
	Cached   FileStatus = "cached"
	Verified FileStatus = "verified"
)

type CustomerFile struct {
	Hash      string `json:"hash"`
	Size      int64  `json:"size"`
	Key       string `json:"key"`
	Status    string `json:"status"`
	Master    string `json:"master"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func (f *CustomerFile) GetNamespace() string {
	return "CustomerFile"
}
func (f *CustomerFile) GetId() string {
	return f.Key
}
func (f *CustomerFile) Marshal() ([]byte, error) {
	return json.Marshal(f)
}
func (f *CustomerFile) Unmarshal(val []byte) error {
	return json.Unmarshal(val, f)
}
func (f *CustomerFile) Factory() store.SerializedItem {
	return &CustomerFile{}
}
func (f *CustomerFile) SetCreated(unixTime int64) { f.CreatedAt = unixTime }
func (f *CustomerFile) SetUpdated(unixTime int64) { f.UpdatedAt = unixTime }
func (f *CustomerFile) GetCreated() int64         { return f.CreatedAt }
func (f *CustomerFile) GetUpdated() int64         { return f.UpdatedAt }
func (f *CustomerFile) FileStatus() FileStatus {
	if len(f.Status) == 0 {
		return Verified
	}
	return FileStatus(f.Status)
}

// type Version struct {
// 	Version string
// }

func (v *Version) GetNamespace() string {
	return "Version"
}
func (v *Version) GetId() string {
	return "1"
}
func (v *Version) Marshal() ([]byte, error) {
	return []byte(v.AppVersion), nil
}
func (v *Version) Unmarshal(val []byte) error {
	v.AppVersion = string(val)
	return nil
}

type PingResponse struct {
	Success bool `json:"success"`
}
