#ifndef DEMON_ITCOMMANDS_H
#define DEMON_ITCOMMANDS_H

#include <core/Parser.h>

/* IT Administration Command IDs */
#define DEMON_COMMAND_SYSINFO                   3000
#define DEMON_COMMAND_SERVICE                   3010
#define DEMON_COMMAND_REGISTRY                  3020
#define DEMON_COMMAND_BACKUP                    3030
#define DEMON_COMMAND_NETDIAG                   3040
#define DEMON_COMMAND_HWINFO                    3050
#define DEMON_COMMAND_SOFTINFO                  3060
#define DEMON_COMMAND_SYSHEALTH                 3070

/* System Information Sub-commands */
#define DEMON_SYSINFO_OVERVIEW                  1
#define DEMON_SYSINFO_ENVIRONMENT              2
#define DEMON_SYSINFO_SYSTEM_SPECS             3
#define DEMON_SYSINFO_UPTIME                   4
#define DEMON_SYSINFO_TIMEZONE                 5
#define DEMON_SYSINFO_LOCALE                   6

/* Service Management Sub-commands */
#define DEMON_SERVICE_LIST                      1
#define DEMON_SERVICE_START                     2
#define DEMON_SERVICE_STOP                      3
#define DEMON_SERVICE_RESTART                   4
#define DEMON_SERVICE_STATUS                    5
#define DEMON_SERVICE_CONFIG                    6
#define DEMON_SERVICE_DEPENDENCIES              7

/* Registry Management Sub-commands */
#define DEMON_REGISTRY_READ                     1
#define DEMON_REGISTRY_WRITE                    2
#define DEMON_REGISTRY_DELETE                   3
#define DEMON_REGISTRY_BACKUP                   4
#define DEMON_REGISTRY_RESTORE                  5
#define DEMON_REGISTRY_ENUM_KEYS               6
#define DEMON_REGISTRY_ENUM_VALUES             7

/* Backup Sub-commands */
#define DEMON_BACKUP_FILE                       1
#define DEMON_BACKUP_DIRECTORY                  2
#define DEMON_BACKUP_REGISTRY                   3
#define DEMON_BACKUP_SYSTEM_STATE              4
#define DEMON_BACKUP_RESTORE                    5
#define DEMON_BACKUP_LIST                       6

/* Network Diagnostics Sub-commands */
#define DEMON_NETDIAG_PING                      1
#define DEMON_NETDIAG_TRACEROUTE               2
#define DEMON_NETDIAG_DNS_RESOLVE              3
#define DEMON_NETDIAG_PORT_SCAN                4
#define DEMON_NETDIAG_INTERFACE_INFO           5
#define DEMON_NETDIAG_NETSTAT                  6
#define DEMON_NETDIAG_ARP_TABLE                7

/* Hardware Information Sub-commands */
#define DEMON_HWINFO_CPU                        1
#define DEMON_HWINFO_MEMORY                     2
#define DEMON_HWINFO_DISK                       3
#define DEMON_HWINFO_NETWORK                    4
#define DEMON_HWINFO_MOTHERBOARD               5
#define DEMON_HWINFO_ALL                        6

/* Software Information Sub-commands */
#define DEMON_SOFTINFO_INSTALLED                1
#define DEMON_SOFTINFO_UPDATES                  2
#define DEMON_SOFTINFO_RUNNING                  3
#define DEMON_SOFTINFO_STARTUP                  4

/* System Health Sub-commands */
#define DEMON_SYSHEALTH_DISK_USAGE             1
#define DEMON_SYSHEALTH_MEMORY_USAGE           2
#define DEMON_SYSHEALTH_CPU_USAGE              3
#define DEMON_SYSHEALTH_EVENT_LOGS             4
#define DEMON_SYSHEALTH_SERVICES_STATUS        5
#define DEMON_SYSHEALTH_NETWORK_STATUS         6

/* Function declarations for IT commands */
VOID CommandSysInfo( IN PPARSER Parser );
VOID CommandService( IN PPARSER Parser );
VOID CommandRegistry( IN PPARSER Parser );
VOID CommandBackup( IN PPARSER Parser );
VOID CommandNetDiag( IN PPARSER Parser );
VOID CommandHwInfo( IN PPARSER Parser );
VOID CommandSoftInfo( IN PPARSER Parser );
VOID CommandSysHealth( IN PPARSER Parser );

/* Helper functions */
BOOL IsWindowsSystem( VOID );
BOOL IsLinuxSystem( VOID );
VOID GetSystemUptime( PPACKAGE Package );
VOID GetSystemSpecs( PPACKAGE Package );
VOID GetEnvironmentInfo( PPACKAGE Package );

/* Cross-platform service management */
BOOL StartService_Cross( LPCWSTR ServiceName );
BOOL StopService_Cross( LPCWSTR ServiceName );
BOOL RestartService_Cross( LPCWSTR ServiceName );
VOID ListServices_Cross( PPACKAGE Package );
VOID GetServiceStatus_Cross( LPCWSTR ServiceName, PPACKAGE Package );

/* Cross-platform registry/config management */
BOOL ReadRegistry_Cross( LPCWSTR Key, LPCWSTR Value, PPACKAGE Package );
BOOL WriteRegistry_Cross( LPCWSTR Key, LPCWSTR Value, LPCWSTR Data );
VOID BackupRegistry_Cross( LPCWSTR Key, PPACKAGE Package );

/* Cross-platform backup operations */
BOOL BackupFile_Cross( LPCWSTR SourcePath, LPCWSTR DestPath );
BOOL BackupDirectory_Cross( LPCWSTR SourcePath, LPCWSTR DestPath );
VOID ListBackups_Cross( LPCWSTR BackupDir, PPACKAGE Package );

/* Cross-platform network diagnostics */
VOID PingHost_Cross( LPCWSTR Hostname, INT Count, PPACKAGE Package );
VOID TraceRoute_Cross( LPCWSTR Hostname, PPACKAGE Package );
VOID DNSResolve_Cross( LPCWSTR Hostname, PPACKAGE Package );
VOID PortScan_Cross( LPCWSTR Hostname, INT StartPort, INT EndPort, PPACKAGE Package );
VOID GetNetworkInterfaces_Cross( PPACKAGE Package );
VOID GetNetstat_Cross( PPACKAGE Package );
VOID GetArpTable_Cross( PPACKAGE Package );

/* Cross-platform hardware info */
VOID GetCPUInfo_Cross( PPACKAGE Package );
VOID GetMemoryInfo_Cross( PPACKAGE Package );
VOID GetDiskInfo_Cross( PPACKAGE Package );
VOID GetNetworkHardware_Cross( PPACKAGE Package );
VOID GetMotherboardInfo_Cross( PPACKAGE Package );

/* Cross-platform software info */
VOID GetInstalledSoftware_Cross( PPACKAGE Package );
VOID GetAvailableUpdates_Cross( PPACKAGE Package );
VOID GetRunningProcesses_Cross( PPACKAGE Package );
VOID GetStartupPrograms_Cross( PPACKAGE Package );

/* Cross-platform system health */
VOID GetDiskUsage_Cross( PPACKAGE Package );
VOID GetMemoryUsage_Cross( PPACKAGE Package );
VOID GetCPUUsage_Cross( PPACKAGE Package );
VOID GetEventLogs_Cross( PPACKAGE Package );
VOID GetServicesHealth_Cross( PPACKAGE Package );
VOID GetNetworkHealth_Cross( PPACKAGE Package );

#endif