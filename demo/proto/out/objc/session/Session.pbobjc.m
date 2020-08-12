// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: session/session.proto

// This CPP symbol can be defined to use imports that match up to the framework
// imports needed when using CocoaPods.
#if !defined(GPB_USE_PROTOBUF_FRAMEWORK_IMPORTS)
 #define GPB_USE_PROTOBUF_FRAMEWORK_IMPORTS 0
#endif

#if GPB_USE_PROTOBUF_FRAMEWORK_IMPORTS
 #import <Protobuf/GPBProtocolBuffers_RuntimeSupport.h>
#else
 #import "GPBProtocolBuffers_RuntimeSupport.h"
#endif

#import "session/Session.pbobjc.h"
// @@protoc_insertion_point(imports)

#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wdeprecated-declarations"

#pragma mark - HLWSessionRoot

@implementation HLWSessionRoot

// No extensions in the file and no imports, so no need to generate
// +extensionRegistry.

@end

#pragma mark - HLWSessionRoot_FileDescriptor

static GPBFileDescriptor *HLWSessionRoot_FileDescriptor(void) {
  // This is called by +initialize so there is no need to worry
  // about thread safety of the singleton.
  static GPBFileDescriptor *descriptor = NULL;
  if (!descriptor) {
    GPB_DEBUG_CHECK_RUNTIME_VERSIONS();
    descriptor = [[GPBFileDescriptor alloc] initWithPackage:@"proto"
                                                 objcPrefix:@"HLW"
                                                     syntax:GPBFileSyntaxProto3];
  }
  return descriptor;
}

#pragma mark - HLWLoginRequest

@implementation HLWLoginRequest

@dynamic name;
@dynamic passWord;

typedef struct HLWLoginRequest__storage_ {
  uint32_t _has_storage_[1];
  NSString *name;
  NSString *passWord;
} HLWLoginRequest__storage_;

// This method is threadsafe because it is initially called
// in +initialize for each subclass.
+ (GPBDescriptor *)descriptor {
  static GPBDescriptor *descriptor = nil;
  if (!descriptor) {
    static GPBMessageFieldDescription fields[] = {
      {
        .name = "name",
        .dataTypeSpecific.className = NULL,
        .number = HLWLoginRequest_FieldNumber_Name,
        .hasIndex = 0,
        .offset = (uint32_t)offsetof(HLWLoginRequest__storage_, name),
        .flags = (GPBFieldFlags)(GPBFieldOptional | GPBFieldTextFormatNameCustom),
        .dataType = GPBDataTypeString,
      },
      {
        .name = "passWord",
        .dataTypeSpecific.className = NULL,
        .number = HLWLoginRequest_FieldNumber_PassWord,
        .hasIndex = 1,
        .offset = (uint32_t)offsetof(HLWLoginRequest__storage_, passWord),
        .flags = (GPBFieldFlags)(GPBFieldOptional | GPBFieldTextFormatNameCustom),
        .dataType = GPBDataTypeString,
      },
    };
    GPBDescriptor *localDescriptor =
        [GPBDescriptor allocDescriptorForClass:[HLWLoginRequest class]
                                     rootClass:[HLWSessionRoot class]
                                          file:HLWSessionRoot_FileDescriptor()
                                        fields:fields
                                    fieldCount:(uint32_t)(sizeof(fields) / sizeof(GPBMessageFieldDescription))
                                   storageSize:sizeof(HLWLoginRequest__storage_)
                                         flags:GPBDescriptorInitializationFlag_None];
#if !GPBOBJC_SKIP_MESSAGE_TEXTFORMAT_EXTRAS
    static const char *extraTextFormatInfo =
        "\002\001D\000\002H\000";
    [localDescriptor setupExtraTextInfo:extraTextFormatInfo];
#endif  // !GPBOBJC_SKIP_MESSAGE_TEXTFORMAT_EXTRAS
    NSAssert(descriptor == nil, @"Startup recursed!");
    descriptor = localDescriptor;
  }
  return descriptor;
}

@end

#pragma mark - HLWLoginReply

@implementation HLWLoginReply

@dynamic flag;
@dynamic cookie;
@dynamic lease;
@dynamic interval;

typedef struct HLWLoginReply__storage_ {
  uint32_t _has_storage_[1];
  int32_t flag;
  int32_t lease;
  int32_t interval;
  NSString *cookie;
} HLWLoginReply__storage_;

// This method is threadsafe because it is initially called
// in +initialize for each subclass.
+ (GPBDescriptor *)descriptor {
  static GPBDescriptor *descriptor = nil;
  if (!descriptor) {
    static GPBMessageFieldDescription fields[] = {
      {
        .name = "flag",
        .dataTypeSpecific.className = NULL,
        .number = HLWLoginReply_FieldNumber_Flag,
        .hasIndex = 0,
        .offset = (uint32_t)offsetof(HLWLoginReply__storage_, flag),
        .flags = (GPBFieldFlags)(GPBFieldOptional | GPBFieldTextFormatNameCustom),
        .dataType = GPBDataTypeInt32,
      },
      {
        .name = "cookie",
        .dataTypeSpecific.className = NULL,
        .number = HLWLoginReply_FieldNumber_Cookie,
        .hasIndex = 1,
        .offset = (uint32_t)offsetof(HLWLoginReply__storage_, cookie),
        .flags = (GPBFieldFlags)(GPBFieldOptional | GPBFieldTextFormatNameCustom),
        .dataType = GPBDataTypeString,
      },
      {
        .name = "lease",
        .dataTypeSpecific.className = NULL,
        .number = HLWLoginReply_FieldNumber_Lease,
        .hasIndex = 2,
        .offset = (uint32_t)offsetof(HLWLoginReply__storage_, lease),
        .flags = (GPBFieldFlags)(GPBFieldOptional | GPBFieldTextFormatNameCustom),
        .dataType = GPBDataTypeInt32,
      },
      {
        .name = "interval",
        .dataTypeSpecific.className = NULL,
        .number = HLWLoginReply_FieldNumber_Interval,
        .hasIndex = 3,
        .offset = (uint32_t)offsetof(HLWLoginReply__storage_, interval),
        .flags = (GPBFieldFlags)(GPBFieldOptional | GPBFieldTextFormatNameCustom),
        .dataType = GPBDataTypeInt32,
      },
    };
    GPBDescriptor *localDescriptor =
        [GPBDescriptor allocDescriptorForClass:[HLWLoginReply class]
                                     rootClass:[HLWSessionRoot class]
                                          file:HLWSessionRoot_FileDescriptor()
                                        fields:fields
                                    fieldCount:(uint32_t)(sizeof(fields) / sizeof(GPBMessageFieldDescription))
                                   storageSize:sizeof(HLWLoginReply__storage_)
                                         flags:GPBDescriptorInitializationFlag_None];
#if !GPBOBJC_SKIP_MESSAGE_TEXTFORMAT_EXTRAS
    static const char *extraTextFormatInfo =
        "\004\001D\000\002F\000\003E\000\004H\000";
    [localDescriptor setupExtraTextInfo:extraTextFormatInfo];
#endif  // !GPBOBJC_SKIP_MESSAGE_TEXTFORMAT_EXTRAS
    NSAssert(descriptor == nil, @"Startup recursed!");
    descriptor = localDescriptor;
  }
  return descriptor;
}

@end

#pragma mark - HLWLogoutRequest

@implementation HLWLogoutRequest

@dynamic id_p;

typedef struct HLWLogoutRequest__storage_ {
  uint32_t _has_storage_[1];
  NSString *id_p;
} HLWLogoutRequest__storage_;

// This method is threadsafe because it is initially called
// in +initialize for each subclass.
+ (GPBDescriptor *)descriptor {
  static GPBDescriptor *descriptor = nil;
  if (!descriptor) {
    static GPBMessageFieldDescription fields[] = {
      {
        .name = "id_p",
        .dataTypeSpecific.className = NULL,
        .number = HLWLogoutRequest_FieldNumber_Id_p,
        .hasIndex = 0,
        .offset = (uint32_t)offsetof(HLWLogoutRequest__storage_, id_p),
        .flags = (GPBFieldFlags)(GPBFieldOptional | GPBFieldTextFormatNameCustom),
        .dataType = GPBDataTypeString,
      },
    };
    GPBDescriptor *localDescriptor =
        [GPBDescriptor allocDescriptorForClass:[HLWLogoutRequest class]
                                     rootClass:[HLWSessionRoot class]
                                          file:HLWSessionRoot_FileDescriptor()
                                        fields:fields
                                    fieldCount:(uint32_t)(sizeof(fields) / sizeof(GPBMessageFieldDescription))
                                   storageSize:sizeof(HLWLogoutRequest__storage_)
                                         flags:GPBDescriptorInitializationFlag_None];
#if !GPBOBJC_SKIP_MESSAGE_TEXTFORMAT_EXTRAS
    static const char *extraTextFormatInfo =
        "\001\001\000ID\000";
    [localDescriptor setupExtraTextInfo:extraTextFormatInfo];
#endif  // !GPBOBJC_SKIP_MESSAGE_TEXTFORMAT_EXTRAS
    NSAssert(descriptor == nil, @"Startup recursed!");
    descriptor = localDescriptor;
  }
  return descriptor;
}

@end

#pragma mark - HLWLogoutReply

@implementation HLWLogoutReply

@dynamic message;

typedef struct HLWLogoutReply__storage_ {
  uint32_t _has_storage_[1];
  NSString *message;
} HLWLogoutReply__storage_;

// This method is threadsafe because it is initially called
// in +initialize for each subclass.
+ (GPBDescriptor *)descriptor {
  static GPBDescriptor *descriptor = nil;
  if (!descriptor) {
    static GPBMessageFieldDescription fields[] = {
      {
        .name = "message",
        .dataTypeSpecific.className = NULL,
        .number = HLWLogoutReply_FieldNumber_Message,
        .hasIndex = 0,
        .offset = (uint32_t)offsetof(HLWLogoutReply__storage_, message),
        .flags = (GPBFieldFlags)(GPBFieldOptional | GPBFieldTextFormatNameCustom),
        .dataType = GPBDataTypeString,
      },
    };
    GPBDescriptor *localDescriptor =
        [GPBDescriptor allocDescriptorForClass:[HLWLogoutReply class]
                                     rootClass:[HLWSessionRoot class]
                                          file:HLWSessionRoot_FileDescriptor()
                                        fields:fields
                                    fieldCount:(uint32_t)(sizeof(fields) / sizeof(GPBMessageFieldDescription))
                                   storageSize:sizeof(HLWLogoutReply__storage_)
                                         flags:GPBDescriptorInitializationFlag_None];
#if !GPBOBJC_SKIP_MESSAGE_TEXTFORMAT_EXTRAS
    static const char *extraTextFormatInfo =
        "\001\001G\000";
    [localDescriptor setupExtraTextInfo:extraTextFormatInfo];
#endif  // !GPBOBJC_SKIP_MESSAGE_TEXTFORMAT_EXTRAS
    NSAssert(descriptor == nil, @"Startup recursed!");
    descriptor = localDescriptor;
  }
  return descriptor;
}

@end


#pragma clang diagnostic pop

// @@protoc_insertion_point(global_scope)
